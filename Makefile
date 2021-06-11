BINARY_NAME=ovnmon
BINARY_MOD=./cmd/$(BINARY_NAME)
MODEL_GEN=$(GOPATH)/bin/modelgen
MODEL_PATH=model
BIN_PATH=bin

.PHONY:all
all: $(BIN_PATH)/$(BINARY_NAME)

.PHONY: clean
clean: 
	@rm -rf $(BIN_PATH)

$(MODEL_GEN):
	@go install github.com/ovn-org/libovsdb/cmd/modelgen

$(BIN_PATH)/$(BINARY_NAME): $(MODEL_GEN)
	@export PATH="$${PATH}:$${GOPATH}/bin";  go generate ./...
	@go build -o $@ $(BINARY_MOD)

	

