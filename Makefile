# Makefile setup
.PHONY: test generate mocks controller-gen doc

# Version locked down by go.mod
# See _code generation_ in README.md
CONTROLLER_GEN ?= go run sigs.k8s.io/controller-tools/cmd/controller-gen

test:
	go test ./... -count=1

all: generate mocks test doc

# Generate code
generate:
	$(CONTROLLER_GEN) object paths="./pkg/apis/..."
	$(CONTROLLER_GEN) crd rbac:roleName=manager-role webhook paths="./pkg/apis/..." output:crd:artifacts:config=config/crd/bases
	cp ./config/crd/bases/*nais.io_*.yaml ./charts/templates

doc:
	mkdir -p doc/output/application
	mkdir -p doc/output/naisjob
	mkdir -p doc/output/topic
	mkdir -p doc/output/openapi/nais
	go run cmd/docgen/docgen.go \
		--dir ./pkg/apis/... \
		--group nais.io \
		--kind Application \
		--reference-output doc/output/application/reference.md \
		--example-output doc/output/application/example.md \
		--openapi-output doc/output/openapi/nais \
		--reference-template doc/templates/reference/application.md \
		--example-template doc/templates/example/application.md \
		;
	go run cmd/docgen/docgen.go \
		--dir ./pkg/apis/... \
		--group nais.io \
		--kind Naisjob \
		--reference-output doc/output/naisjob/reference.md \
		--example-output doc/output/naisjob/example.md \
		--openapi-output doc/output/openapi/nais \
		--reference-template doc/templates/reference/naisjob.md \
		--example-template doc/templates/example/naisjob.md \
		;
	go run cmd/docgen/docgen.go \
		--dir ./pkg/apis/... \
		--group kafka.nais.io \
		--kind Topic \
		--reference-output doc/output/topic/reference.md \
		--example-output doc/output/topic/example.md \
		--openapi-output doc/output/openapi/nais \
		--reference-template doc/templates/reference/topic.md \
		--example-template doc/templates/example/topic.md \
		;

mocks:
	go run github.com/vektra/mockery/v2
