package cliconfig

import (
	"fmt"
	"github.com/ContainX/depcon/utils"
	"github.com/bgentry/speakeasy"
	"net/url"
	"os"
	"regexp"
)

const (
	AlphaNumDash string = `^[a-zA-Z0-9_-]*$`
)

var (
	rxAlphaNumDash *regexp.Regexp = regexp.MustCompile(AlphaNumDash)
)

func getAlpaNumDash(question string) string {
	var response string

	fmt.Printf("%s: ", question)
	fmt.Scanf("%s", &response)

	if rxAlphaNumDash.MatchString(response) {
		return response
	}

	fmt.Printf("\nERROR: '%s' must contain valid characters within %s\n", response, AlphaNumDash)
	return getAlpaNumDash(question)
}

func getPassword(question string) string {
	password, err := speakeasy.Ask(fmt.Sprintf("%s", question))
	if err != nil {
		fmt.Printf("\nERROR: %s\n", err.Error())
		return getPassword(question)
	}
	return password
}

func getPasswordWithVerify() string {
	pass1 := getPassword("Password: ")
	pass2 := getPassword("Verify Password: ")
	if pass1 != pass2 {
		fmt.Println("Password and Verify Password don't match\n")
		return getPasswordWithVerify()
	}
	return pass1
}

// Asks the user for the remote URI of the Marathon service
func getMarathonURL(count int) string {
	if count > 5 {
		fmt.Printf("Too many retries obtaining Marathon URL.  If depcon is running within docker please insure 'docker run -it' is set.\n")
		os.Exit(1)
	}
	var response string
	fmt.Print("Marathon URL (eg. http://hostname:8080)  : ")
	fmt.Scanf("%s", &response)

	_, err := url.ParseRequestURI(response)
	if err == nil && utils.HasURLScheme(response) {
		return response
	}

	fmt.Printf("\nERROR: '%s' must be a valid URL\n", response)
	return getMarathonURL(count + 1)
}

func createEnvironment() *ServiceConfig {
	service := ServiceConfig{}
	service.Name = getAlpaNumDash("Environment Name (eg. test, stage, prod) ")
	service.HostUrl = getMarathonURL(0)

	if getBoolAnswer("Authentication Required", false) {
		service.Username = getAlpaNumDash("Username")
		service.Password = getPasswordWithVerify()
	}
	fmt.Println("")
	return &service
}
