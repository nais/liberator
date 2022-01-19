package kafka_nais_io_v1

import (
	"testing"

	aiven_nais_io_v1 "github.com/nais/liberator/pkg/apis/aiven.nais.io/v1"
	"github.com/stretchr/testify/assert"
)

func TestAclNameFromTopicAcl(t *testing.T) {
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
		{
			name: "test max length",
			args: args{
				acl: &TopicACL{
					Access:      "read",
					Application: "team-superlong-team-name-a-very-long-application-name-that-needs-to-be-shortened",
					Team:        "team-superlong-team-name",
				},
				suffix: "99",
			},
			want: "superlong-team-name_a-very-long-application-name-t_aef7fe79_99",
		},
		{
			name: "wildcards",
			args: args{
				acl: &TopicACL{
					Access:      "write",
					Application: "*",
					Team:        "*",
				},
				suffix: "00",
			},
			want: "*_*_*_00",
		},
		{
			name: "wildcard app",
			args: args{
				acl: &TopicACL{
					Access:      "write",
					Application: "*",
					Team:        "myteam",
				},
				suffix: "00",
			},
			want: "myteam_*_*_00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.acl.ServiceUserNameWithSuffix(tt.args.suffix)
			assert.NoError(t, err, "AclNameFromTopicAcl(%v, %v)", tt.args.acl, tt.args.suffix)
			assert.Equalf(t, tt.want, got, "AclNameFromTopicAcl(%v, %v)", tt.args.acl, tt.args.suffix)
			assert.LessOrEqual(t, len(got), aiven_nais_io_v1.MaxServiceUserNameLength, "length of service username too long")
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
		{
			name: "avoid separator at end",
			args: args{
				team:        "myteam",
				application: "myteam-a-long-application-name-which-needs-to-be-shortened",
			},
			want: "a-long-application-name-which",
		},
		{
			name: "wildcard appname",
			args: args{
				team:        "myteam",
				application: "*",
			},
			want: "*",
		},
		{
			name: "wildcard team and app",
			args: args{
				team:        "*",
				application: "*",
			},
			want: "*",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, shortAppName(tt.args.team, tt.args.application), "shortAppName(%v, %v)", tt.args.team, tt.args.application)
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
		{
			name: "avoid separator at end",
			args: args{"team-superlong-team-name-actually-very-long"},
			want: "superlong-team-name",
		},
		{
			name: "wildcard",
			args: args{"*"},
			want: "*",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, shortTeamName(tt.args.team), "shortTeamName(%v)", tt.args.team)
		})
	}
}
