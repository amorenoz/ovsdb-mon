BINARY_NAME=ovnmon
BINARY_MOD=./cmd/$(BINARY_NAME)
MODEL_GEN=$(GOPATH)/bin/modelgen
MODEL_PATH=model
BIN_PATH=bin
SCHEMA?=schemas/ovn-nb.ovsschema

.PHONY:all
all: build

$(MODEL_GEN):
ifeq ($(GOPATH),)
	$(error GOPATH is not set)
endif
	@go install github.com/ovn-org/libovsdb/cmd/modelgen

.PHONY: build
build: $(MODEL_GEN)
	@echo "Generating model based on schema $(SCHEMA)"
	@cp $(SCHEMA) model/db.ovsschema
	@export PATH="$${PATH}:$${GOPATH}/bin";  go generate ./...
	@go build -o $(BIN_PATH)/$(BINARY_NAME) $(BINARY_MOD)

.PHONY: clean
clean: 
	@rm -rf $(BIN_PATH)
	@find model -name "*.go" -not -name "gen.go" | xargs rm -f
	@rm -f model/db.ovsschema

