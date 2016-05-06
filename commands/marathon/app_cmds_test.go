package marathon

import (
	l "log"
	"testing"
)

func TestParseParamFile(t *testing.T) {
	envParams, _ := parseParamsFile("resources/test.env")
	el, ok := envParams["APP1_VERSION"]
	if !ok && el != "3" {
		l.Printf("Actual envParams %v", envParams)
		l.Panic("Expected envFile parsed correctly")
	}
	el, ok = envParams["APP2_VERSION"]
	if !ok && el != "345" {
		l.Printf("Actual envParams %v", envParams)
		l.Panic("Expected envFile parsed correctly")
	}
}
