package kafka_nais_io_v1_test

import (
	"sort"
	"testing"

	"github.com/nais/liberator/pkg/apis/kafka.nais.io/v1"
	"github.com/stretchr/testify/assert"
)

type UserList []kafka_nais_io_v1.User

func (ul UserList) Len() int {
	return len(ul)
}

func (ul UserList) Less(i, j int) bool {
	return ul[i].Username < ul[j].Username
}

func (ul UserList) Swap(i, j int) {
	x := ul[i]
	ul[i] = ul[j]
	ul[j] = x
}

func TestTopicACLs_Usernames(t *testing.T) {
	acls := kafka_nais_io_v1.TopicACLs{
		{
			Application: "app",
			Team:        "team",
		},
		{ // duplicate
			Application: "app",
			Team:        "team",
		},
		{
			Application: "app2",
			Team:        "team",
		},
		{
			Application: "app3",
			Team:        "team2",
		},
	}

	expected := UserList{
		{
			Username:    "team.app-407e3d92",
			Application: "app",
			Team:        "team",
		},
		{
			Username:    "team.app2-4943258",
			Application: "app2",
			Team:        "team",
		},
		{
			Username:    "team2.app3-8c3bd36a",
			Application: "app3",
			Team:        "team2",
		},
	}

	actual := UserList(acls.Users())

	sort.Sort(expected)
	sort.Sort(actual)

	assert.Equal(t, expected, actual)
}
