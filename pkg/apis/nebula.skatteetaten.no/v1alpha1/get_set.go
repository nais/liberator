package v1alpha1

import (
	"fmt"

	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	nais_io_v1alpha1 "github.com/nais/liberator/pkg/apis/nais.io/v1alpha1"
)

// TODO: replace manual getters with generated code
// TODO: candidates are either `go generate` or a switch to Protobuf

func (in *Application) SetStatus(status *nais_io_v1.Status) {
	in.Status = *status
}

func (in *Application) SetStatusConditions() {
	in.Status.SetStatusConditions()
}

func (in *Application) GetStatus() *nais_io_v1.Status {
	return &in.Status
}

func (in *Application) GetStrategy() *nais_io_v1.Strategy {
	return &nais_io_v1.Strategy{
		Type: nais_io_v1alpha1.DeploymentStrategyRollingUpdate,
	}
}

func (in *Application) GetReplicas() *nais_io_v1.Replicas {
	return in.Spec.Replicas
}

func (in *Application) GetPrometheus() *nais_io_v1.PrometheusConfig {
	return &nais_io_v1.PrometheusConfig{
		Enabled: !in.Spec.Pod.Prometheus.Disabled,
		Port:    in.Spec.Pod.Prometheus.Port,
		Path:    in.Spec.Pod.Prometheus.Path,
	}
}

func (in *Application) GetLogtransform() string {
	return ""
}

func (in *Application) GetLogformat() string {
	return ""
}

func (in *Application) GetPort() int {
	return nais_io_v1alpha1.DefaultAppPort
}

func (in *Application) GetService() *nais_io_v1.Service {
	return &nais_io_v1.Service{
		Port:     nais_io_v1alpha1.DefaultServicePort,
		Protocol: nais_io_v1alpha1.DefaultPortName,
	}
}

func (in *Application) GetLiveness() *nais_io_v1.Probe {
	return in.Spec.Pod.Liveness
}

func (in *Application) GetReadiness() *nais_io_v1.Probe {
	return in.Spec.Pod.Readiness
}

func (in *Application) GetStartup() *nais_io_v1.Probe {
	return in.Spec.Pod.Startup
}

func (in *Application) GetEnv() nais_io_v1.EnvVars {
	return in.Spec.Pod.Env
}

func (in *Application) GetImage() string {
	return fmt.Sprintf("%s:%s", in.Spec.Pod.ImageName, in.Spec.Pod.Tag)
}

func (in *Application) GetImageName() string {
	return in.Spec.Pod.ImageName
}

func (in *Application) GetCommand() []string {
	return in.Spec.Pod.Command
}

func (in *Application) GetFilesFrom() []nais_io_v1.FilesFrom {
	return in.Spec.Pod.FilesFrom
}

func (in *Application) GetEnvFrom() []nais_io_v1.EnvFrom {
	return in.Spec.Pod.EnvFrom
}

func (in *Application) GetPreStopHook() *nais_io_v1.PreStopHook {
	return nil
}

func (in *Application) GetPreStopHookPath() string {
	return ""
}

func (in *Application) GetResources() *nais_io_v1.ResourceRequirements {
	return in.Spec.Pod.Resources
}

func (in *Application) IsAuroraApplication() bool {
	return in.Spec.AuroraApplication == true
}

// GetLogDirectoryOr will find the configured log directory, or use the
// indicated default log dir
func (in *Application) GetLogDirectoryOr(defaultLogDir string) string {

	if in.Spec.Logging == nil || in.Spec.Logging.LogDirectory == "" {
		return defaultLogDir
	} else {
		return in.Spec.Logging.LogDirectory
	}
}

func (in *Application) GetIngress() *IngressConfig {
	publicIngress := map[string]PublicIngressConfig{}
	for key, item := range in.Spec.Ingress.Public {
		if item.Disabled {
			continue
		}
		// TODO: it is hard to get kubebuilder default values applied in test so we do it here.
		if item.Port == 0 {
			item.Port = 8080
		}

		if item.ServicePort == 0 {
			item.ServicePort = 80
		}

		if item.Gateway == "" {
			item.Gateway = "istio-ingressgateway"
		}
		publicIngress[key] = item
	}

	items := map[string]InternalIngressConfig{}
	for key, item := range in.Spec.Ingress.Internal {
		if item.Disabled {
			continue
		}

		if item.Namespace == "" {
			item.Namespace = in.Namespace
		}

		items[key] = item
	}
	return &IngressConfig{
		Public:   publicIngress,
		Internal: items,
	}
}

func (in *Application) GetEgress() *EgressConfig {

	if in.Spec.Egress == nil {
		return &EgressConfig{}
	}
	externalEgress := map[string]ExternalEgressConfig{}
	for key, item := range in.Spec.Egress.External {
		if item.Disabled {
			continue
		}

		if len(item.Ports) == 0 {
			item.Ports = []PortConfig{{
				Port:     443,
				Protocol: "HTTPS",
				Name:     "https",
			}}
		}
		externalEgress[key] = item
	}

	items := map[string]InternalEgressConfig{}
	for key, item := range in.Spec.Egress.Internal {
		if item.Disabled {
			continue
		}
		if item.Namespace == "" {
			item.Namespace = in.Namespace
		}
		items[key] = item
	}
	return &EgressConfig{
		External: externalEgress,
		Internal: items,
	}
}

func (in *Application) GetImagePolicy() *ImagePolicyConfig {
	if in.Spec.ImagePolicy == nil || in.Spec.ImagePolicy.Disabled {
		return nil
	}
	return in.Spec.ImagePolicy
}

func (in *Application) GetAzure() *AzureConfig {
	return in.Spec.Azure
}

func (in *Application) GetAzureResourceGroup() string {
	return in.Spec.Azure.ResourceGroup
}

func (in *Application) GetPostgresDatabases() map[string]*PostgreDatabaseConfig {
	dbs := map[string]*PostgreDatabaseConfig{}

	// TODO: how do we validate that the config is correct here, if you have multiple with primary it will fail
	// Is this where we would need to use a validation hook?
	for key, db := range in.Spec.Azure.PostgreDatabases {
		if db.Disabled {
			continue
		}
		users := map[string]*PostgreDatabaseUser{}
		for uk, user := range db.Users {
			if user.Disabled {
				continue
			}

			if user.Prefix == "" {
				user.Prefix = "SPRING_DATASOURCE"
			}
			users[uk] = user
		}
		dbs[key] = db
	}

	// if we only have a single database and a single user mark it as primary by default
	if len(dbs) == 1 {
		for key, db := range dbs {
			if len(db.Users) == 1 {
				for user, _ := range db.Users {
					dbs[key].Users[user].Primary = true
				}
			}
		}
	}
	return dbs
}

func (in *Application) GetStorageAccounts() map[string]*StorageAccountConfig {
	items := map[string]*StorageAccountConfig{}
	for key, item := range in.Spec.Azure.StorageAccount {
		if item.Disabled {
			continue
		}

		if item.Prefix == "" {
			item.Prefix = "AZURE_STORAGE"
		}
		items[key] = item
	}

	// if single item set it as primary
	if len(items) == 1 {
		for key, _ := range items {
			items[key].Primary = true
		}
	}

	return items
}

func (in *Application) GetCosmosDb() map[string]*CosmosDBConfig {
	items := map[string]*CosmosDBConfig{}
	for key, item := range in.Spec.Azure.CosmosDB {
		if item.Disabled {
			continue
		}

		if item.Prefix == "" {
			if item.MongoDBVersion != "" {
				item.Prefix = "SPRING_DATA_MONGODB"
			} else {
				item.Prefix = "COSMOSDB"
			}
		}
		items[key] = item
	}

	// if single item set it as primary
	if len(items) == 1 {
		for key, _ := range items {
			items[key].Primary = true
		}
	}
	return items
}
