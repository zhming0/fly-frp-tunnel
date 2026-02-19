IMG ?= fly-frp-tunnel:latest

.PHONY: all
all: build

## Build
.PHONY: build
build:
	go build -o bin/manager .

.PHONY: run
run:
	go run .

.PHONY: test
test:
	go test ./... -v

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

## Docker
.PHONY: docker-build
docker-build:
	docker build -t $(IMG) .

.PHONY: docker-push
docker-push:
	docker push $(IMG)

## Install/Deploy (using Helm)
.PHONY: helm-install
helm-install:
	helm install fly-frp-tunnel charts/fly-frp-tunnel

.PHONY: helm-uninstall
helm-uninstall:
	helm uninstall fly-frp-tunnel

.PHONY: helm-template
helm-template:
	helm template fly-frp-tunnel charts/fly-frp-tunnel

## Manifests
.PHONY: manifests
manifests:
	@echo "No CRDs to generate — this operator uses only core Kubernetes types."

.PHONY: clean
clean:
	rm -rf bin/
