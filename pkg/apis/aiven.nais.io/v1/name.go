package aiven_nais_io_v1

func ValkeyFullyQualifiedName(instance, namespace string) string {
	return "valkey-" + namespace + "-" + instance
}

func OpenSearchFullyQualifiedName(instance, namespace string) string {
	return "opensearch-" + namespace + "-" + instance
}
