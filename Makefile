BINARY_NAME=ovnmon
BINARY_MOD=./cmd/$(BINARY_NAME)
GEN_NAME=modelgen
GEN_MOD=./cmd/$(GEN_NAME)
MODEL_PATH=model
BIN_PATH=bin

.PHONY:all
all: $(BIN_PATH)/$(BINARY_NAME)

.PHONY: clean
clean: 
	@rm -r $(BIN_PATH)
	@rm -r $(MODEL_PATH)


$(BIN_PATH)/$(GEN_NAME): $(GEN_MOD)
	go build -o $@ ./$^

$(MODEL_PATH): $(BIN_PATH)/$(GEN_NAME)
	$(BIN_PATH)/$(GEN_NAME) -o $@ -p model ovn-nb.ovsschema

$(BIN_PATH)/$(BINARY_NAME): $(MODEL_PATH)
	go build -o $@ $(BINARY_MOD)

	

