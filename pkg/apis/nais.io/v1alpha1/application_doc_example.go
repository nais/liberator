package nais_io_v1alpha1

import (
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ExampleApplicationForDocumentation() *Application {
	intp := func(i int) *int {
		return &i
	}

	return &Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: "nais.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "myapplication",
			Namespace: "myteam",
			Labels: map[string]string{
				"team": "myteam",
			},
		},
		Spec: ApplicationSpec{
			AccessPolicy: &nais_io_v1.AccessPolicy{
				Inbound: &nais_io_v1.AccessPolicyInbound{
					Rules: []nais_io_v1.AccessPolicyRule{
						{
							Application: "app1",
						},
						{
							Application: "app2",
							Namespace:   "q1",
						},
						{
							Application: "app3",
							Namespace:   "q2",
							Cluster:     "dev-gcp",
						},
						{
							Application: "*",
							Namespace:   "q3",
						},
					},
				},
				Outbound: &nais_io_v1.AccessPolicyOutbound{
					Rules: []nais_io_v1.AccessPolicyRule{
						{
							Application: "app1",
						},
						{
							Application: "app2",
							Namespace:   "q1",
						},
						{
							Application: "app3",
							Namespace:   "q2",
							Cluster:     "dev-gcp",
						},
						{
							Application: "*",
							Namespace:   "q3",
						},
					},
					External: []nais_io_v1.AccessPolicyExternalRule{
						{
							Host: "external-application.example.com",
						},
						{
							Host: "non-http-service.example.com",
							Ports: []nais_io_v1.AccessPolicyPortRule{
								{
									Name:     "kafka",
									Port:     9200,
									Protocol: "TCP",
								},
							},
						},
					},
				},
			},
			Azure: &nais_io_v1.Azure{
				Application: &nais_io_v1.AzureApplication{
					Enabled: true,
					ReplyURLs: []string{
						"https://myapplication.nav.no/oauth2/callback",
					},
					Tenant: "nav.no",
					Claims: &nais_io_v1.AzureAdClaims{
						Extra: []nais_io_v1.AzureAdExtraClaim{
							"NAVident",
							"azp_name",
						},
						Groups: []nais_io_v1.AzureAdGroup{
							{
								ID: "00000000-0000-0000-0000-000000000000",
							},
						},
					},
				},
			},
			Command: []string{
				"/app/myapplication",
				"--param",
				"value",
				"--other-param",
				"other-value",
			},
			Elastic: &nais_io_v1.Elastic{
				Instance: "my-elastic-instance",
			},
			Env: []nais_io_v1.EnvVar{
				{
					Name:  "MY_CUSTOM_VAR",
					Value: "some_value",
				},
				{
					Name: "MY_APPLICATION_NAME",
					ValueFrom: &nais_io_v1.EnvVarSource{
						FieldRef: nais_io_v1.ObjectFieldSelector{
							FieldPath: "metadata.name",
						},
					},
				},
			},
			EnvFrom: []nais_io_v1.EnvFrom{
				{
					Secret: "my-secret-with-envs",
				},
				{
					ConfigMap: "my-configmap-with-envs",
				},
			},
			FilesFrom: []nais_io_v1.FilesFrom{
				{
					ConfigMap: "example-files-configmap",
					MountPath: "/var/run/configmaps",
				},
				{
					Secret:    "my-secret-file",
					MountPath: "/var/run/secrets",
				},
			},
			GCP: &nais_io_v1.GCP{
				Buckets: []nais_io_v1.CloudStorageBucket{
					{
						Name:                "my-cloud-storage-bucket",
						CascadingDelete:     true,
						RetentionPeriodDays: intp(30),
						LifecycleCondition: &nais_io_v1.LifecycleCondition{
							Age:              10,
							CreatedBefore:    "2020-01-01",
							NumNewerVersions: 2,
							WithState:        "ARCHIVED",
						},
					},
				},
				SqlInstances: []nais_io_v1.CloudSqlInstance{
					{
						Type:             "POSTGRES_12",
						Name:             "myinstance",
						Tier:             "db-f1-micro",
						DiskType:         "SSD",
						HighAvailability: true,
						DiskSize:         30,
						DiskAutoresize:   true,
						AutoBackupHour:   intp(1),
						Maintenance: &nais_io_v1.Maintenance{
							Day:  1,
							Hour: intp(4),
						},
						Databases: []nais_io_v1.CloudSqlDatabase{
							{
								Name:         "mydatabase",
								EnvVarPrefix: "DB",
								Users: []nais_io_v1.CloudSqlDatabaseUser{
									{
										Name: "extra_user",
									},
								},
							},
						},
						CascadingDelete: true,
						Collation:       "nb_NO.UTF8",
					},
				},
				Permissions: []nais_io_v1.CloudIAMPermission{
					{
						Role: "roles/cloudsql.client",
						Resource: nais_io_v1.CloudIAMResource{
							APIVersion: "resourcemanager.cnrm.cloud.google.com/v1beta1",
							Kind:       "Project",
							Name:       "myteam-dev-ab23",
						},
					},
				},
			},
			IDPorten: &nais_io_v1.IDPorten{
				AccessTokenLifetime:    intp(3600),
				ClientURI:              "https://www.nav.no",
				Enabled:                true,
				FrontchannelLogoutPath: "/oauth2/logout",
				FrontchannelLogoutURI:  "https://myapplication.nav.no/oauth2/logout",
				PostLogoutRedirectURIs: []string{
					"https://www.nav.no",
				},
				RedirectPath:    "/oauth2/callback",
				RedirectURI:     "https://myapplication.nav.no/oauth2/callback",
				SessionLifetime: intp(7200),
			},
			Image: "navikt/testapp:69.0.0",
			Ingresses: []nais_io_v1.Ingress{
				"https://myapplication.nav.no",
			},
			Kafka: &nais_io_v1.Kafka{
				Pool: "nav-dev",
			},
			LeaderElection: true,
			Liveness: &nais_io_v1.Probe{
				FailureThreshold: 10,
				InitialDelay:     20,
				Path:             "/isalive",
				PeriodSeconds:    5,
				Port:             8080,
				Timeout:          1,
			},
			Logformat:    "accesslog_with_referer_useragent",
			Logtransform: "http_loglevel",
			Maskinporten: &nais_io_v1.Maskinporten{
				Enabled: true,
				Scopes: []nais_io_v1.MaskinportenScope{
					{
						Name: "some_scope",
					},
				},
			},
			Port:            8080,
			PreStopHookPath: "/internal/stop",
			Prometheus: &nais_io_v1.PrometheusConfig{
				Enabled: true,
				Port:    "8080",
				Path:    "/metrics",
			},
			Readiness: &nais_io_v1.Probe{
				FailureThreshold: 10,
				InitialDelay:     20,
				Path:             "/isready",
				PeriodSeconds:    5,
				Port:             8080,
				Timeout:          1,
			},
			Replicas: &nais_io_v1.Replicas{
				Min:                    2,
				Max:                    4,
				CpuThresholdPercentage: 50,
			},
			Resources: &nais_io_v1.ResourceRequirements{
				Limits: &nais_io_v1.ResourceSpec{
					Cpu:    "500m",
					Memory: "512Mi",
				},
				Requests: &nais_io_v1.ResourceSpec{
					Cpu:    "200m",
					Memory: "256Mi",
				},
			},
			SecureLogs: &nais_io_v1.SecureLogs{
				Enabled: true,
			},
			Service: &nais_io_v1.Service{
				Port:     DefaultServicePort,
				Protocol: DefaultPortName,
			},
			SkipCaBundle: true,
			Startup: &nais_io_v1.Probe{
				FailureThreshold: 10,
				InitialDelay:     20,
				Path:             "/started",
				PeriodSeconds:    5,
				Port:             8080,
				Timeout:          1,
			},
			Strategy: &nais_io_v1.Strategy{
				Type: DeploymentStrategyRollingUpdate,
			},
			TokenX: &nais_io_v1.TokenX{
				Enabled:                 true,
				MountSecretsAsFilesOnly: true,
			},
			Vault: &nais_io_v1.Vault{
				Enabled: true,
				Sidecar: true,
				Paths: []nais_io_v1.SecretPath{
					{
						MountPath: "/var/run/secrets/nais.io/vault",
						KvPath:    "/kv/preprod/fss/application/namespace",
						Format:    "env",
					},
				},
			},
			WebProxy: true,
		},
	}
}
