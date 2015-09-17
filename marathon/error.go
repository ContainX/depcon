package marathon

import (
	"errors"
)

var (
	ErrorTimeout            = errors.New("The operation has timed out")
	ErrorDeploymentNotfound = errors.New("Failed to get deployment in allocated time")
)
