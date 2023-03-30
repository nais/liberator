package nais_io_v1alpha1

import (
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"github.com/nais/liberator/pkg/intutil"
)

func ExampleApplicationForDocumentation() *Application {
	intp := func(i int) *int {
		return &i
	}
	int64p := func(i int64) *int64 {
		return &i
	}
	boolp := func(b bool) *bool {
		return &b
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
					Rules: []nais_io_v1.AccessPolicyInboundRule{
						{
							AccessPolicyRule: nais_io_v1.AccessPolicyRule{
								Application: "app1",
							},
						},
						{
							AccessPolicyRule: nais_io_v1.AccessPolicyRule{
								Application: "app2",
								Namespace:   "q1",
							},
						},
						{
							AccessPolicyRule: nais_io_v1.AccessPolicyRule{
								Application: "app3",
								Namespace:   "q2",
								Cluster:     "dev-gcp",
							},
						},
						{
							AccessPolicyRule: nais_io_v1.AccessPolicyRule{
								Application: "*",
								Namespace:   "q3",
							},
						},
						{
							AccessPolicyRule: nais_io_v1.AccessPolicyRule{
								Application: "app4",
							},
							Permissions: &nais_io_v1.AccessPolicyPermissions{
								Scopes: []nais_io_v1.AccessPolicyPermission{"custom-scope"},
							},
						},
						{
							AccessPolicyRule: nais_io_v1.AccessPolicyRule{
								Application: "app5",
							},
							Permissions: &nais_io_v1.AccessPolicyPermissions{
								Roles: []nais_io_v1.AccessPolicyPermission{"custom-role"},
							},
						},
						{
							AccessPolicyRule: nais_io_v1.AccessPolicyRule{
								Application: "app6",
							},
							Permissions: &nais_io_v1.AccessPolicyPermissions{
								Scopes: []nais_io_v1.AccessPolicyPermission{"custom-scope"},
								Roles:  []nais_io_v1.AccessPolicyPermission{"custom-role"},
							},
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
					Enabled:       true,
					AllowAllUsers: boolp(true),
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
					ReplyURLs: []nais_io_v1.AzureAdReplyUrlString{
						"https://myapplication.nav.no/oauth2/callback",
					},
					SinglePageApplication: boolp(true),
					Tenant:                "nav.no",
				},
				Sidecar: &nais_io_v1.AzureSidecar{
					Wonderwall: nais_io_v1.Wonderwall{
						AutoLogin: true,
						AutoLoginIgnorePaths: []nais_io_v1.WonderwallIgnorePaths{
							"/path",
							"/internal/*",
						},
						Enabled:   true,
						ErrorPath: "/error",
						Resources: &nais_io_v1.ResourceRequirements{
							Limits: &nais_io_v1.ResourceSpec{
								Cpu:    "250m",
								Memory: "256Mi",
							},
							Requests: &nais_io_v1.ResourceSpec{
								Cpu:    "20m",
								Memory: "32Mi",
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
				{
					EmptyDir: &nais_io_v1.EmptyDir{
						Medium: "Memory",
					},
					MountPath: "/var/cache",
				},
				{
					PersistentVolumeClaim: "pvc-name",
					MountPath:             "/var/run/pvc",
				},
			},
			GCP: &nais_io_v1.GCP{
				BigQueryDatasets: []nais_io_v1.CloudBigQueryDataset{
					{
						Name:            "my_bigquery_dataset1",
						CascadingDelete: true,
						Description:     "Contains big data, supporting big queries, for use in big ideas.",
						Permission:      nais_io_v1.BigQueryPermissionReadWrite,
					},
					{
						Name:        "my_bigquery_dataset2",
						Description: "Contains big data, supporting big queries, for use in big ideas.",
						Permission:  nais_io_v1.BigQueryPermissionRead,
					},
				},
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
						UniformBucketLevelAccess: true,
						PublicAccessPrevention:   true,
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
						RetainedBackups:  intp(14),
						Maintenance: &nais_io_v1.Maintenance{
							Day:  1,
							Hour: intp(4),
						},
						Flags: []nais_io_v1.CloudSqlFlag{
							{
								Name:  "max_connections",
								Value: "50",
							},
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
						CascadingDelete:     true,
						Collation:           "nb_NO.UTF8",
						PointInTimeRecovery: true,
						Insights: &nais_io_v1.InsightsConfiguration{
							Enabled:               boolp(true),
							QueryStringLength:     4500,
							RecordApplicationTags: true,
							RecordClientAddress:   true,
						},
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
				IntegrationType:        "idporten",
				PostLogoutRedirectURIs: []nais_io_v1.IDPortenURI{
					"https://www.nav.no",
				},
				RedirectPath:    "/oauth2/callback",
				Scopes:          []string{"openid", "profile"},
				SessionLifetime: intp(7200),
				Sidecar: &nais_io_v1.IDPortenSidecar{
					Wonderwall: nais_io_v1.Wonderwall{
						AutoLogin: true,
						AutoLoginIgnorePaths: []nais_io_v1.WonderwallIgnorePaths{
							"/path",
							"/internal/*",
						},
						Enabled:   true,
						ErrorPath: "/error",
						Resources: &nais_io_v1.ResourceRequirements{
							Limits: &nais_io_v1.ResourceSpec{
								Cpu:    "250m",
								Memory: "256Mi",
							},
							Requests: &nais_io_v1.ResourceSpec{
								Cpu:    "20m",
								Memory: "32Mi",
							},
						},
					},
					Level:  "Level4",
					Locale: "nb",
				},
			},
			Influx: &nais_io_v1.Influx{
				Instance: "influx-instance",
			},
			Image: "navikt/testapp:69.0.0",
			Ingresses: []nais_io_v1.Ingress{
				"https://myapplication.nav.no",
			},
			Kafka: &nais_io_v1.Kafka{
				Pool:    "nav-dev",
				Streams: true,
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
				Scopes: nais_io_v1.MaskinportenScope{
					ConsumedScopes: []nais_io_v1.ConsumedScope{
						{
							Name: "skatt:scope.read",
						},
					},
					ExposedScopes: []nais_io_v1.ExposedScope{
						{
							Enabled:             true,
							Name:                "scope.read",
							Product:             "arbeid",
							AllowedIntegrations: []string{"maskinporten"},
							AtMaxAge:            intp(30),
							Consumers: []nais_io_v1.ExposedScopeConsumer{
								{
									Orgno: "123456789",
									Name:  "KST",
								},
							},
						},
					},
				},
			},
			OpenSearch: &nais_io_v1.OpenSearch{
				Instance: "my-open-search-instance",
				Access:   "readwrite",
			},
			Port: 8080,
			PreStopHook: &nais_io_v1.PreStopHook{
				Exec: &nais_io_v1.ExecAction{
					Command: []string{"./my", "--shell", "script"},
				},
				Http: &nais_io_v1.HttpGetAction{
					Path: "/internal/stop",
					Port: intp(8080),
				},
			},
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
				Min:                    intutil.Intp(2),
				Max:                    intutil.Intp(4),
				CpuThresholdPercentage: 50,
				DisableAutoScaling:     true,
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
				RollingUpdate: &v1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(0),
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.String,
						StrVal: "25%",
					},
				},
			},
			TerminationGracePeriodSeconds: int64p(60),
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
