VERSION?=$(shell git describe --tags --always)
DOCKER_IMAGE=quintilesims/k8-auth0-authenticator

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --ldflags "-X main.Version=$(VERSION)" -o k8-auth0-authenticator . 
	docker build -t $(DOCKER_IMAGE):$(VERSION) .

release: build
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker tag  $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest
	docker push $(DOCKER_IMAGE):latest

.PHONY: build release
