# Makefile setup
.PHONY: test generate mocks controller-gen

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
	$(CONTROLLER_GEN) crd:trivialVersions=true rbac:roleName=manager-role webhook paths="./pkg/apis/nais.io/... ./pkg/apis/kafka.nais.io/..." output:crd:artifacts:config=config/crd/bases

mocks:
	cd pkg/ && mockery -inpkg -all -case snake

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@$$(CONTROLLER_GEN_VERSION) ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
