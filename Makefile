.PHONY: build build-dev push deploy generate wire lint test

REGISTRY ?= k3d-muxly-registry.localhost:5111
TAG ?= dev
IMAGE := $(REGISTRY)/muxly-msg-subscriber:$(TAG)

build:
	docker build --build-arg GITHUB_TOKEN=$$GITHUB_PERSONAL_ACCESS_TOKEN --target prod -t $(IMAGE) .

build-dev:
	docker build --build-arg GITHUB_TOKEN=$$GITHUB_PERSONAL_ACCESS_TOKEN --target dev -t $(IMAGE) .

push:
	docker push $(IMAGE)

deploy: build push
	kubectl rollout restart deployment/muxly-msg-subscriber -n muxly

generate:
	./oapi-codegen.sh

wire:
	wire gen ./internal/app/

lint:
	golangci-lint run ./...

test:
	go test ./...
