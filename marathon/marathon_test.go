package marathon

import "testing"

func TestNewMarathonClient(t *testing.T) {
	opts := &MarathonOptions{}
	client := createMarathonClient("username", "password", "", opts, "localhost")

	if client == nil {
		t.Error()
	}

	var i interface{} = client
	_, ok := i.(*MarathonClient)
	if !ok {
		t.Error()
	}
}
