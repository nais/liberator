# Makefile setup
.PHONY: test generate mocks controller-gen doc

# Version locked down by go.mod
# See _code generation_ in README.md
CONTROLLER_GEN ?= go run sigs.k8s.io/controller-tools/cmd/controller-gen

test:
	go test ./... -count=1

# Generate code
generate:
	$(CONTROLLER_GEN) object paths="./pkg/apis/..."
	$(CONTROLLER_GEN) crd:trivialVersions=true,preserveUnknownFields=false rbac:roleName=manager-role webhook paths="./pkg/apis/..." output:crd:artifacts:config=config/crd/bases
	cp ./config/crd/bases/*nais.io_*.yaml ./charts/templates

doc:
	mkdir -p doc/output/application
	mkdir -p doc/output/naisjob
	mkdir -p doc/output/alert
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
		--group nais.io \
		--kind Alert \
		--reference-output doc/output/alert/reference.md \
		--example-output doc/output/alert/example.md \
		--openapi-output doc/output/openapi/nais \
		--reference-template doc/templates/reference/alert.md \
		--example-template doc/templates/example/alert.md \
		;

mocks:
	cd pkg/aiven/ && mockery --inpackage --all --case snake
