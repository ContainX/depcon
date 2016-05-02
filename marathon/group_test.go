package marathon

import (
	"github.com/ContainX/depcon/pkg/mockrest"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	GroupsFolder = "testdata/groups/"
)

func TestListGroups(t *testing.T) {
	s := mockrest.StartNewWithFile(GroupsFolder + "list_groups_response.json")
	defer s.Stop()

	c := NewMarathonClient(s.URL, "", "")
	groups, err := c.ListGroups()

	assert.Nil(t, err, "Error response was not expected")
	assert.Equal(t, 1, len(groups.Groups), "Expected 1 nested group")
	assert.Equal(t, 0, len(groups.Apps), "Expected 0 top level apps")
}

func TestGetGroup(t *testing.T) {
	s := mockrest.StartNewWithFile(GroupsFolder + "get_group_response.json")
	defer s.Stop()

	c := NewMarathonClient(s.URL, "", "")
	group, err := c.GetGroup("/sites")

	assert.Nil(t, err, "Error response was not expected")
	assert.Equal(t, "/sites", group.GroupID, "Expected /sites identifier")
}

func TestDestroyGroup(t *testing.T) {
	s := mockrest.StartNewWithFile(CommonFolder + "deployid_response.json")
	defer s.Stop()

	c := NewMarathonClient(s.URL, "", "")
	depId, err := c.DestroyGroup("/sites")
	assert.Nil(t, err, "Error response was not expected")
	assert.Equal(t, "5ed4c0c5-9ff8-4a6f-a0cd-f57f59a34b43", depId.DeploymentID)
}
