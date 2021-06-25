package nais_io_v1_test

import (
	"testing"

	"github.com/nais/liberator/pkg/apis/nais.io/v1"

	"github.com/stretchr/testify/assert"
)

const secretName = "verysecret"

var accessPolicy = &nais_io_v1.AccessPolicy{
	Inbound: &nais_io_v1.AccessPolicyInbound{
		Rules: []nais_io_v1.AccessPolicyInboundRule{
			{
				AccessPolicyRule: nais_io_v1.AccessPolicyRule{
					Application: "app1",
					Namespace:   "ns1",
					Cluster:     "firstcluster",
				},
			},
		},
	},
	Outbound: &nais_io_v1.AccessPolicyOutbound{
		Rules: []nais_io_v1.AccessPolicyRule{
			{
				Application: "app1",
				Namespace:   "ns1",
				Cluster:     "firstcluster",
			},
		},
	},
}

func TestHash(t *testing.T) {
	spec := nais_io_v1.JwkerSpec{
		AccessPolicy: accessPolicy,
		SecretName:   secretName,
	}
	hash, err := spec.Hash()
	assert.NoError(t, err)
	assert.Equal(t, "b6b03694476d8028", hash)
}
