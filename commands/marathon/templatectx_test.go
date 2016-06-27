package marathon

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeFunctionality(t *testing.T) {
	tc, err := LoadTemplateContext("resources/testcontext.json")
	assert.NoError(t, err)

	m := tc.mergeAppWithDefault("prod")

	assert.Equal(t, float64(0.1), m["appa"]["cpus"])
	assert.Equal(t, float64(3), m["appa"]["instances"])
	assert.Equal(t, float64(300), m["appa"]["mem"])

}
