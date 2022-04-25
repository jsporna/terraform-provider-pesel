VERSION ?= 0.1.0

NAME = pesel
BINARY = terraform-provider-${NAME}
MARCH = "$$(go env GOOS)_$$(go env GOARCH)"

ACCTEST_PARALLELISM ?= 10
ACCTEST_TIMEOUT = 120m
ACCTEST_COUNT = 1
TEST ?= ./...

export GOBIN = $(shell pwd)/bin

$(GOBIN):
	mkdir -p $(GOBIN)

.PHONY: build
build: lint
	go build -o ${BINARY}

.PHONY: testacc
testacc: lint
	TF_ACC=1 go test -v ./... -count $(ACCTEST_COUNT) -parallel $(ACCTEST_PARALLELISM) -timeout $(ACCTEST_TIMEOUT)

.PHONY: test
test: lint
	go test -v $(TEST) $(TESTARGS) -timeout=5m -parallel=4

.PHONY: docs-generate
docs-generate: tools
	@ $(GOBIN)/tfplugindocs

.PHONY: gen
gen: docs-generate
	@ go generate ./...

.PHONY: clean
clean:
	rm -f ${BINARY}

.PHONY: install
install: build ## Install built provider into the local terraform cache
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/jsporna/${NAME}/${VERSION}/${MARCH}
	mv ${BINARY} ~/.terraform.d/plugins/registry.terraform.io/jsporna/${NAME}/${VERSION}/${MARCH}

.PHONY: tools
tools: $(GOBIN)
	@ cd tools && go install github.com/client9/misspell/cmd/misspell
	@ cd tools && go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	@ cd tools && go install github.com/golangci/golangci-lint/cmd/golangci-lint
	@ cd tools && go install github.com/goreleaser/goreleaser

.PHONY: misspell
misspell:
	@ $(GOBIN)/misspell -error -source go ./internal/
	@ $(GOBIN)/misspell -error -source text ./templates/

.PHONY: golangci-lint
golangci-lint:
	@ $(GOBIN)/golangci-lint run --max-same-issues=0 --timeout=300s $(GOLANGCIFLAGS) ./internal/...

.PHONY: lint
lint: setup misspell golangci-lint

.PHONY: setup
setup: tools
