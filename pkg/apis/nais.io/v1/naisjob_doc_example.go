package nais_io_v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ExampleNaisjobForDocumentation() *Naisjob {
	intp := func(i int) *int {
		return &i
	}
	int32p := func(i int32) *int32 {
		return &i
	}
	int64p := func(i int64) *int64 {
		return &i
	}

	return &Naisjob{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Naisjob",
			APIVersion: "nais.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "myjob",
			Namespace: "myteam",
			Labels: map[string]string{
				"team": "myteam",
			},
		},
		Spec: NaisjobSpec{
			AccessPolicy: &AccessPolicy{
				Inbound: &AccessPolicyInbound{
					Rules: []AccessPolicyRule{
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
				Outbound: &AccessPolicyOutbound{
					Rules: []AccessPolicyRule{
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
					External: []AccessPolicyExternalRule{
						{
							Host: "external-application.example.com",
						},
						{
							Host: "non-http-service.example.com",
							Ports: []AccessPolicyPortRule{
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
			ActiveDeadlineSeconds: int64p(60),
			Azure: &Azure{
				Application: &AzureApplication{
					Enabled: true,
					ReplyURLs: []string{
						"https://myapplication.nav.no/oauth2/callback",
					},
					Tenant: "nav.no",
					Claims: &AzureAdClaims{
						Extra: []AzureAdExtraClaim{
							"NAVident",
							"azp_name",
						},
						Groups: []AzureAdGroup{
							{
								ID: "00000000-0000-0000-0000-000000000000",
							},
						},
					},
				},
			},
			BackoffLimit: 5,
			Command: []string{
				"/app/myapplication",
				"--param",
				"value",
				"--other-param",
				"other-value",
			},
			Elastic: &Elastic{
				Instance: "my-elastic-instance",
			},
			Env: []EnvVar{
				{
					Name:  "MY_CUSTOM_VAR",
					Value: "some_value",
				},
				{
					Name: "MY_APPLICATION_NAME",
					ValueFrom: &EnvVarSource{
						FieldRef: ObjectFieldSelector{
							FieldPath: "metadata.name",
						},
					},
				},
			},
			EnvFrom: []EnvFrom{
				{
					Secret: "my-secret-with-envs",
				},
				{
					ConfigMap: "my-configmap-with-envs",
				},
			},
			FailedJobsHistoryLimit: 2,
			FilesFrom: []FilesFrom{
				{
					ConfigMap: "example-files-configmap",
					MountPath: "/var/run/configmaps",
				},
				{
					Secret:    "my-secret-file",
					MountPath: "/var/run/secrets",
				},
			},
			GCP: &GCP{
				BigQueryDatasets: []CloudBigQueryDataset{
					{
						Name:            "my_bigquery_dataset1",
						CascadingDelete: true,
						Description:     "Contains big data, supporting big queries, for use in big ideas.",
						Permission:      BigQueryPermissionReadWrite,
					},
					{
						Name:            "my_bigquery_dataset2",
						Description:     "Contains big data, supporting big queries, for use in big ideas.",
						Permission:      BigQueryPermissionRead,
					},
				},
				Buckets: []CloudStorageBucket{
					{
						Name:                "my-cloud-storage-bucket",
						CascadingDelete:     true,
						RetentionPeriodDays: intp(30),
						LifecycleCondition: &LifecycleCondition{
							Age:              10,
							CreatedBefore:    "2020-01-01",
							NumNewerVersions: 2,
							WithState:        "ARCHIVED",
						},
					},
				},
				SqlInstances: []CloudSqlInstance{
					{
						Type:             "POSTGRES_12",
						Name:             "myinstance",
						Tier:             "db-f1-micro",
						DiskType:         "SSD",
						HighAvailability: true,
						DiskSize:         30,
						DiskAutoresize:   true,
						AutoBackupHour:   intp(1),
						Maintenance: &Maintenance{
							Day:  1,
							Hour: intp(4),
						},
						Databases: []CloudSqlDatabase{
							{
								Name:         "mydatabase",
								EnvVarPrefix: "DB",
								Users: []CloudSqlDatabaseUser{
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
				Permissions: []CloudIAMPermission{
					{
						Role: "roles/cloudsql.client",
						Resource: CloudIAMResource{
							APIVersion: "resourcemanager.cnrm.cloud.google.com/v1beta1",
							Kind:       "Project",
							Name:       "myteam-dev-ab23",
						},
					},
				},
			},
			Influx: &Influx{
				Instance: "influx-instance",
			},
			Image: "navikt/testapp:69.0.0",
			Kafka: &Kafka{
				Pool: "nav-dev",
			},
			Liveness: &Probe{
				FailureThreshold: 10,
				InitialDelay:     20,
				Path:             "/isalive",
				PeriodSeconds:    5,
				Port:             8080,
				Timeout:          1,
			},
			Logformat:    "accesslog_with_referer_useragent",
			Logtransform: "http_loglevel",
			Maskinporten: &Maskinporten{
				Enabled: true,
				Scopes: MaskinportenScope{
					ConsumedScopes: []ConsumedScope{
						{
							Name: "skatt:scope.read",
						},
					},
					ExposedScopes: []ExposedScope{
						{
							Enabled:             true,
							Name:                "scope.read",
							Product:             "arbeid",
							AllowedIntegrations: []string{"maskinporten"},
							Consumers: []ExposedScopeConsumer{
								{
									Orgno: "123456789",
									Name:  "KST",
								},
							},
						},
					},
				},
			},
			PreStopHook: &PreStopHook{
				Exec: &ExecAction{
					Command: []string{"./my", "--shell", "script"},
				},
				Http: &HttpGetAction{
					Path: "/internal/stop",
					Port: intp(8080),
				},
			},
			PreStopHookPath: "/internal/stop",
			Readiness: &Probe{
				FailureThreshold: 10,
				InitialDelay:     20,
				Path:             "/isready",
				PeriodSeconds:    5,
				Port:             8080,
				Timeout:          1,
			},
			Resources: &ResourceRequirements{
				Limits: &ResourceSpec{
					Cpu:    "500m",
					Memory: "512Mi",
				},
				Requests: &ResourceSpec{
					Cpu:    "200m",
					Memory: "256Mi",
				},
			},
			Schedule: "*/15 0 0 0 0",
			SecureLogs: &SecureLogs{
				Enabled: true,
			},
			SkipCaBundle: true,
			Startup: &Probe{
				FailureThreshold: 10,
				InitialDelay:     20,
				Path:             "/started",
				PeriodSeconds:    5,
				Port:             8080,
				Timeout:          1,
			},
			SuccessfulJobsHistoryLimit: 2,
			TTLSecondsAfterFinished:    int32p(60),
			Vault: &Vault{
				Enabled: true,
				Sidecar: true,
				Paths: []SecretPath{
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
