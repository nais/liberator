# Makefile setup
.PHONY: test generate mocks controller-gen doc

# Lock down version of controller-gen
# See _code generation_ in README.md
CONTROLLER_GEN_VERSION ?= "v0.2.5"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

test:
	go test ./... -count=1

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object paths="./pkg/apis/..."
	$(CONTROLLER_GEN) crd:trivialVersions=true rbac:roleName=manager-role webhook paths="./pkg/apis/..." output:crd:artifacts:config=config/crd/bases

doc:
	mkdir -p doc/output/application
	mkdir -p doc/output/naisjob
	go run cmd/docgen/docgen.go \
		--dir ./pkg/apis/... \
		--group nais.io \
		--kind Application \
		--reference-output doc/output/application/reference.md \
		--example-output doc/output/application/example.md \
		--reference-template doc/templates/reference/application.md \
		--example-template doc/templates/example/application.md \
		;
	go run cmd/docgen/docgen.go \
		--dir ./pkg/apis/... \
		--group nais.io \
		--kind Naisjob \
		--reference-output doc/output/naisjob/reference.md \
		--example-output doc/output/naisjob/example.md \
		--reference-template doc/templates/reference/naisjob.md \
		--example-template doc/templates/example/naisjob.md \
		;

mocks:
	cd pkg/ && mockery --inpackage --all --case snake

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION) ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
