package kafka_nais_io_v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAclName(t *testing.T) {
	type args struct {
		acl    *TopicACL
		suffix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test redundant information",
			args: args{
				acl: &TopicACL{
					Access:      "read",
					Application: "team-redundant-team-application",
					Team:        "team-redundant-team",
				},
				suffix: "*",
			},
			want: "redundant-team_application_18515795_*",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AclName(tt.args.acl, tt.args.suffix)
			assert.NoError(t, err, "AclName(%v, %v)", tt.args.acl, tt.args.suffix)
			assert.Equalf(t, tt.want, got, "AclName(%v, %v)", tt.args.acl, tt.args.suffix)
		})
	}
}

func Test_shortAppName(t *testing.T) {
	type args struct {
		team        string
		application string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no shortening",
			args: args{
				team:        "myteam",
				application: "myapp",
			},
			want: "myapp",
		},
		{
			name: "drop team prefix",
			args: args{
				team:        "myteam",
				application: "myteammyapp",
			},
			want: "myapp",
		},
		{
			name: "drop team prefix with separator",
			args: args{
				team:        "myteam",
				application: "myteam-myapp",
			},
			want: "myapp",
		},
		{
			name: "cut long name",
			args: args{
				team:        "myteam",
				application: "a-very-long-application-name-that-needs-to-be-shortened",
			},
			want: "a-very-long-application-name-t",
		},
		{
			name: "cut long name with team prefix",
			args: args{
				team:        "myteam",
				application: "myteam-a-very-long-application-name-that-needs-to-be-shortened",
			},
			want: "a-very-long-application-name-t",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acl := &TopicACL{
				Application: tt.args.application,
				Team:        tt.args.team,
			}
			assert.Equalf(t, tt.want, shortAppName(acl), "shortAppName(%v)", acl)
		})
	}
}

func Test_shortTeamName(t *testing.T) {
	type args struct {
		team string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no shortening",
			args: args{"myteam"},
			want: "myteam",
		},
		{
			name: "drop prefix",
			args: args{"teammyteam"},
			want: "myteam",
		},
		{
			name: "drop prefix and separator",
			args: args{"team-myteam"},
			want: "myteam",
		},
		{
			name: "shorten long name",
			args: args{"a-very-long-team-name-that-needs-to-be-shortened"},
			want: "a-very-long-team-nam",
		},
		{
			name: "shorten long name with prefix",
			args: args{"team-a-very-long-team-name-that-needs-to-be-shortened"},
			want: "a-very-long-team-nam",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, shortTeamName(tt.args.team), "shortTeamName(%v)", tt.args.team)
		})
	}
}
