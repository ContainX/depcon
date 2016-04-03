package bluegreen

import (
	"encoding/csv"
	"fmt"
	"github.com/gondor/depcon/marathon"
	"github.com/gondor/depcon/utils"
	"io"
	"math"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	HAProxyStatsQP = "/haproxy?stats;csv"
	HAProxyPidsQP  = "/_haproxy_getpids"
	BackendRE      = `(?i)^(\d+)_(\d+)_(\d+)_(\d+)_(\d+)$`
)

// Simple HTTP Test to determine if the current LB is the correct URL. Better to test this before we modify Marathon with this
// existing deployment
func (c *BGClient) isProxyAlive() {
	resp := c.http.HttpGet(c.opts.LoadBalancer+HAProxyStatsQP, nil)
	if resp.Error != nil {
		log.Fatal("HAProxy is not responding or is invalid or Stats service not enabled.\n\n", resp.Error.Error())
	}
}

func (c *BGClient) checkIfTasksDrained(app, existingApp *marathon.Application, stepStartedAt time.Time) bool {
	time.Sleep(c.opts.StepDelay)

	existingApp = c.refreshAppOrPanic(existingApp.ID)
	app = c.refreshAppOrPanic(app.ID)

	targetInstances, _ := strconv.Atoi(app.Labels[DeployTargetInstances])
	log.Info("Existing app running %d instance, new app running %d instances", existingApp.Instances, app.Instances)

	hosts, err := proxiesFromURI(c.opts.LoadBalancer)
	if err != nil {
		log.Error("Error with HAProxy Stats URL: %s", err.Error())
	}

	var errCaught error = nil
	var csvData string

	for _, h := range hosts {
		log.Debug("Querying HAProxy stats: %s", h+HAProxyStatsQP)
		resp := c.http.HttpGet(h+HAProxyStatsQP, nil)
		if resp.Error != nil {
			errCaught = resp.Error
		} else {
			csvData = csvData + resp.Content
		}

		resp = c.http.HttpGet(h+HAProxyPidsQP, nil)
		if resp.Error != nil {
			errCaught = resp.Error
		} else {
			pids := strings.Split(resp.Content, " ")
			if len(pids) > 1 && time.Now().Sub(stepStartedAt) < c.opts.StepDelay {
				log.Info("Waiting for %d, pids on %s", len(pids), h)
				return c.checkIfTasksDrained(app, existingApp, stepStartedAt)
			}
		}

		if errCaught != nil {
			log.Warning("Caught error when retrieving HAProxy stats from %s: Error (%s)", h, errCaught.Error())
			return c.checkIfTasksDrained(app, existingApp, stepStartedAt)
		}
	}

	pinfo := parseProxyBackends(csvData, app)
	if len(pinfo.backends)/pinfo.instanceCount != (app.Instances + existingApp.Instances) {
		// HAProxy hasn't updated yet, try again
		return c.checkIfTasksDrained(app, existingApp, stepStartedAt)
	}

	backendsUp := backendsForStatus(pinfo, "UP")
	if len(backendsUp)/pinfo.instanceCount < targetInstances {
		// Wait until we're in a health state
		return c.checkIfTasksDrained(app, existingApp, stepStartedAt)
	}

	// Double check that current draining backends are finished serving requests
	backendsDrained := backendsForStatus(pinfo, "MAINT")
	if len(backendsDrained)/pinfo.instanceCount < 1 {
		// No backends have started draining yet
		return c.checkIfTasksDrained(app, existingApp, stepStartedAt)
	}

	for _, be := range backendsDrained {
		// Verify that the backends have no sessions or pending connections.
		// This is likely overkill, but we'll do it anyway to be safe.
		if intOrZero(string(be[pinfo.hmap["qcur"]])) > 0 || intOrZero(string(be[pinfo.hmap["scur"]])) > 0 {
			// Backends are not yet defined
			return c.checkIfTasksDrained(app, existingApp, stepStartedAt)
		}
	}

	// If we made it here, all the backends are drained and we can start removing tasks, with prejudice
	hostPorts := hostPortsFromBackends(pinfo.hmap, backendsDrained, pinfo.instanceCount)
	tasksToKill := findTasksToKill(existingApp.Tasks, hostPorts)

	log.Info("There are %d drained backends, about to kill & scale for these tasks:\n%s", len(tasksToKill), strings.Join(tasksToKill, "\n"))

	if app.Instances == targetInstances && len(tasksToKill) == existingApp.Instances {
		log.Info("About to delete old app %s", existingApp.ID)
		if _, err := c.marathon.DestroyApplication(existingApp.ID); err != nil {
			return false
		}
		return true
	}

	// Scale new app up
	instances := int(math.Floor(float64(app.Instances + (app.Instances+1)/2)))
	if instances >= existingApp.Instances {
		instances = targetInstances
	}
	log.Info("Scaling new app up to %d instances", instances)
	if _, err := c.marathon.ScaleApplication(app.ID, instances); err != nil {
		panic("Failed to scale application: " + err.Error())
	}

	//Scale old app down
	log.Info("Scaling old app down to %d instances", len(tasksToKill))
	if err := c.marathon.KillTasksAndScale(tasksToKill...); err != nil {
		log.Error("Failure killing tasks: %v", tasksToKill)
	}
	return c.checkIfTasksDrained(app, existingApp, time.Now())

}

func findTasksToKill(tasks []*marathon.Task, hostPorts map[string][]int) []string {
	tasksToKill := map[string]string{}
	for _, task := range tasks {
		if _, ok := hostPorts[task.Host]; ok {
			for _, p := range hostPorts[task.Host] {
				if utils.IntInSlice(p, task.Ports) {
					tasksToKill[task.ID] = task.ID
				}
			}
		}
	}
	return utils.MapStringKeysToSlice(tasksToKill)
}

func hostPortsFromBackends(hmap map[string]int, backends [][]string, instanceCount int) map[string][]int {
	regex := regexp.MustCompile(BackendRE)
	counts := map[string]int{}
	hostPorts := map[string][]int{}

	for _, be := range backends {
		svname := string(be[hmap["svname"]])
		if _, ok := counts[svname]; ok {
			counts[svname] += 1
		} else {
			counts[svname] = 1
		}

		if counts[svname] == instanceCount {
			if regex.MatchString(svname) {
				m := regex.FindStringSubmatch(svname)
				host := strings.Join(m[1:5], ".")
				port := m[5]

				if _, ok := hostPorts[host]; ok {
					hostPorts[host] = append(hostPorts[host], intOrZero(port))
				} else {
					hostPorts[host] = []int{intOrZero(port)}
				}
			}
		}
	}
	return hostPorts
}

func backendsForStatus(pinfo *proxyInfo, status string) [][]string {
	var results [][]string
	for _, b := range pinfo.backends {
		if b[pinfo.hmap["status"]] == status {
			results = append(results, b)
		}
	}
	return results
}

func parseProxyBackends(data string, app *marathon.Application) *proxyInfo {

	pi := &proxyInfo{
		instanceCount: 0,
		hmap:          map[string]int{},
		backends:      make([][]string, 0),
	}

	var headers []string

	r := csv.NewReader(strings.NewReader(data))
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if []rune(row[0])[0] == '#' {
			headers = row
			pi.instanceCount += 1
			continue
		}

		if row[0] == fmt.Sprintf("%s_%s", app.Labels[DeployGroup], app.Labels[DeployProxyPort]) && notBackFrontend(row[1]) {
			pi.backends = append(pi.backends, row)
		}
	}

	log.Info("Found %d backends across %d HAProxy instances", len(pi.backends), pi.instanceCount)

	// Create header map of column to index
	for i := 0; i < len(headers); i++ {
		pi.hmap[headers[i]] = i
	}
	return pi
}

func notBackFrontend(value string) bool {
	return value != "BACKEND" && value != "FRONTEND"
}

func (c *BGClient) refreshAppOrPanic(id string) *marathon.Application {
	// Retry in case of minor network errors
	for i := 0; i < 3; i++ {
		if a, err := c.marathon.GetApplication(id); err != nil {
			log.Error("Error refresh app info: %s, Will retry %d more times before giving up", err.Error(), 3-(i+1))
			time.Sleep(time.Duration(3) * time.Second)
		} else {
			return a
		}
	}
	panic("Failure to refresh application " + id)
}

func proxiesFromURI(uri string) ([]string, error) {
	url, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	ips, err := net.LookupIP(url.Host)
	if err != nil {
		return []string{url.String()}, nil
	}

	results := []string{}
	for _, ip := range ips {
		url.Host = ip.String()
		results = append(results, url.String())
	}
	return results, nil
}
