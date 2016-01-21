package marathon

import (
	"testing"
	"log"
)

func TestParseParamFile(t *testing.T) {
	envParams, _ := parseParamsFile("resources/test.env")
	el, ok := envParams["APP1_VERSION"]
	if !ok && el != "3" {
		log.Printf("Actual envParams %v", envParams)
		log.Panic("Expected envFile parsed correctly")
	}
	el, ok = envParams["APP2_VERSION"]
	if !ok && el != "345" {
		log.Printf("Actual envParams %v", envParams)
		log.Panic("Expected envFile parsed correctly")
	}
}
