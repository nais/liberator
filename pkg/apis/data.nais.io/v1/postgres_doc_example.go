package v1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func ExamplePostgresForDocumentation() *Postgres {
	return &Postgres{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Postgres",
			APIVersion: "data.nais.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mycluster",
			Namespace: "myteam",
			Labels: map[string]string{
				"team": "myteam",
			},
		},
		Spec: PostgresSpec{
			Cluster: PostgresCluster{
				Resources: PostgresResources{
					DiskSize: resource.MustParse("2Gi"),
					Cpu:      resource.MustParse("200m"),
					Memory:   resource.MustParse("2Gi"),
				},
				MajorVersion:     "17",
				HighAvailability: true,
				AllowDeletion:    true,
				Audit: &PostgresAudit{
					Enabled:          true,
					StatementClasses: []PostgresAuditStatementClass{"function", "misc"},
				},
			},
			Database: &PostgresDatabase{
				Collation: "nb_NO.UTF-8",
				Extensions: []PostgresExtension{
					{
						Name: "postgis",
					},
				},
			},
			MaintenanceWindow: &Maintenance{
				Day:  4,
				Hour: ptr.To(10),
			},
		},
	}
}
