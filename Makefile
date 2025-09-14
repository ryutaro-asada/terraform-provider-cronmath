HOSTNAME=registry.terraform.io
NAMESPACE=ryutaro-asada
NAME=cronmath
BINARY=terraform-provider-${NAME}
VERSION=1.0.0

# Determine OS and Architecture
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)

default: install

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test -v -cover ./...

testacc:
	TF_ACC=1 go test -v -timeout 120m ./...

docs:
	go generate ./...

fmt:
	go fmt ./...
	gofmt -s -w .

lint:
	golangci-lint run ./...

clean:
	rm -f ${BINARY}
	rm -rf ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}
	find . -name ".terraform" -type d -exec rm -rf {} +
	find . -name ".terraform.lock.hcl" -type f -delete
	find . -name "terraform.tfstate*" -type f -delete

release:
	goreleaser release --clean

snapshot:
	goreleaser release --snapshot --clean

.PHONY: build install test testacc docs fmt lint clean release snapshot
