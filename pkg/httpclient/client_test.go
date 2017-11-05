package httpclient

import (
	"sync"
	"testing"
)

func TestHttpClient_Configuration(t *testing.T) {
	httpConfig := NewDefaultConfig()
	httpConfig.HttpUser = "username"
	httpConfig.HttpPass = "passwd"
	httpConfig.HttpToken = ""
	httpConfig.RWMutex = sync.RWMutex{}

	client := NewHttpClient(httpConfig)
	if client == nil {
		t.Error()
	}
}
