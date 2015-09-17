package cliconfig

import (
	"github.com/gondor/depcon/utils"
	"fmt"
	"github.com/bgentry/speakeasy"
	"net/url"
	"regexp"
)

const (
	AlphaNumDash string = `^[a-zA-Z0-9_]*$`
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
func getMarathonURL() string {
	var response string
	fmt.Print("Marathon URL (eg. http://hostname:8080)  : ")
	fmt.Scanf("%s", &response)

	_, err := url.ParseRequestURI(response)
	if err == nil && utils.HasURLScheme(response) {
		return response
	}

	fmt.Printf("\nERROR: '%s' must be a valid URL\n", response)
	return getMarathonURL()
}

func createEnvironment() *ServiceConfig {
	service := ServiceConfig{}
	service.Name = getAlpaNumDash("Environment Name (eg. test, stage, prod) ")
	service.HostUrl = getMarathonURL()

	if getBoolAnswer("Authentication Required", false) {
		service.Username = getAlpaNumDash("Username")
		service.Password = getPasswordWithVerify()
	}
	fmt.Println("")
	return &service
}
