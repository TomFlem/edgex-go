#
# Copyright (c) 2018 Cavium
#
# SPDX-License-Identifier: Apache-2.0
#


.PHONY: build clean test docker run


GO=CGO_ENABLED=0 go
GOCGO=CGO_ENABLED=1 go

DOCKERS=docker_export_client docker_export_distro docker_core_data docker_core_metadata docker_core_command docker_support_logging
.PHONY: $(DOCKERS)

MICROSERVICES=cmd/export-client/export-client cmd/export-distro/export-distro cmd/core-metadata/core-metadata cmd/core-data/core-data cmd/core-command/core-command cmd/support-logging/support-logging
.PHONY: $(MICROSERVICES)

VERSION=$(shell cat ./VERSION)

GOFLAGS=-ldflags "-X github.com/edgexfoundry/edgex-go.Version=$(VERSION)"

build: $(MICROSERVICES)
	go build ./...

cmd/core-metadata/core-metadata:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/core-metadata

cmd/core-data/core-data:
	$(GOCGO) build $(GOFLAGS) -o $@ ./cmd/core-data

cmd/core-command/core-command:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/core-command

cmd/export-client/export-client:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/export-client

cmd/export-distro/export-distro:
	$(GOCGO) build $(GOFLAGS) -o $@ ./cmd/export-distro

cmd/support-logging/support-logging:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/support-logging

clean:
	rm -f $(MICROSERVICES)

test:
	go test ./...
	go vet ./...

prepare:
	glide install

run:
	cd bin && ./edgex-launch.sh

run_docker:
	cd bin && ./edgex-docker-launch.sh

docker: $(DOCKERS)

docker_core_metadata:
	docker build -f docker/Dockerfile.core-metadata -t edgexfoundry/core-metadata-go:$(VERSION) .

docker_core_data:
	docker build -f docker/Dockerfile.core-data -t edgexfoundry/core-data-go:$(VERSION) .

docker_core_command:
	docker build -f docker/Dockerfile.core-command -t edgexfoundry/core-command-go:$(VERSION) .

docker_export_client:
	docker build -f docker/Dockerfile.export-client -t edgexfoundry/export-client-go:$(VERSION) .

docker_export_distro:
	docker build -f docker/Dockerfile.export-distro -t edgexfoundry/export-distro-go:$(VERSION) .

docker_support_logging:
	docker build -f docker/Dockerfile.support-logging -t edgexfoundry/support-logging-go:$(VERSION) .
