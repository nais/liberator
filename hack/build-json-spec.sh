#!/bin/bash

set -e

K8S_VERSION=v1.19.12

OPENAPI_DIR=$(pwd)/doc/output/openapi

# Heavily inspired by https://github.com/yannh/kubernetes-json-schema/blob/master/build.sh
OPENAPI2JSONSCHEMABIN="docker create -i ghcr.io/yannh/openapi2jsonschema:latest"

# Generate json spec for k8s resources
container=$($OPENAPI2JSONSCHEMABIN -o "schemas/kubernetes" --kubernetes --stand-alone --expanded --strict "https://raw.githubusercontent.com/kubernetes/kubernetes/${K8S_VERSION}/api/openapi-spec/swagger.json")
set +e
docker start -ai "$container"
set -e
rm -rf "$OPENAPI_DIR/kubernetes"
docker cp "$container:/out/schemas/kubernetes" "$OPENAPI_DIR"
docker rm "$container"

# Make a json file with all nais resources
echo '{"oneOf":[' > "$OPENAPI_DIR/nais-all.json"
list=$(ls "$OPENAPI_DIR/nais/")
mapfile -t FILES <<< "$list"
for file in ${FILES[*]}; do
	echo "{\"\$ref\":\"nais/$(basename "$file")\"}," >> "$OPENAPI_DIR/nais-all.json"
done
truncate -s-2 "$OPENAPI_DIR/nais-all.json"
echo "]}" >> "$OPENAPI_DIR/nais-all.json"

# Combine nais and k8s resources, but only those with an enum defined for the kind (otherwise there's completion problems)
jq_k8s='
	{
		"oneOf":
		(
			.[1].oneOf +
				(
					.[0].definitions
					| to_entries
					| map(select(.value.properties.kind.enum))
					| map({"$ref": ("kubernetes/_definitions.json#/definitions/"+ .key)})
				)
		)
	}
'
jq -s "$jq_k8s" "$OPENAPI_DIR/kubernetes/_definitions.json" "$OPENAPI_DIR/nais-all.json" > \
		"$OPENAPI_DIR/nais-k8s-all.json"

# Upload to k8s
# gsutil -m  -h "Cache-Control:private, max-age=0, no-transform" cp -r "$OPENAPI_DIR/"{kubernetes,nais,nais-all.json,nais-k8s-all.json} $BUCKET
