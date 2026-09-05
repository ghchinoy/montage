BINARY  := montage
BIN_DIR := bin
PKG     := .

.PHONY: build build-ui build-embedded clean test vet fmt install run all docs-diagrams deploy

all: build-embedded

deploy:
	./scripts/deploy.sh

build-ui:
	cd ui && npm run build

build-embedded: build-ui
	@mkdir -p $(BIN_DIR)
	go build -tags embedui -o $(BIN_DIR)/$(BINARY) $(PKG)

docs-diagrams:
	dot -Twebp docs/architecture.dot -o docs/architecture.webp
	dot -Twebp docs/pipeline.dot -o docs/pipeline.webp

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY) $(PKG)

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

install:
	go install .

clean:
	rm -rf $(BIN_DIR)
	rm -f $(BINARY)

run: build
	./$(BIN_DIR)/$(BINARY)
