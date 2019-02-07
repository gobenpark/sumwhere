GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
IMAGE=qkrqjadn/sumwhere
DOCKER_PASSWORD=1q2w3e4r
DOCKER_USERNAME=qkrqjadn
VERSION=1.0
GITCOMMITCOUNT:=$$(git rev-list HEAD | wc -l | tr -d ' ')
GITHASH:=$$(git rev-parse --short HEAD)
DATETIME:=$$(date "+%Y%m%d-%H%M%S")
VERSIONS:=$(VERSION).$(GITCOMMITCOUNT)-$(GITHASH)-$(DATETIME)
#https://codecov.io/

.PHONY: clean build-docker rolling-update sumwhere

clean:
	$(GOCLEAN)

sumwhere: clean
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o $@ -ldflags "-X main.ServiceVersion=$(VERSIONS)" *.go

build-docker: build-alpine
	@docker build -t $(IMAGE):$(TAG) .
	@docker push $(IMAGE):$(TAG)

rolling-update: build-docker
	@ssh root@202.30.23.76 -p 55555 kubectl set image deployment/sumwhere-server sumwhere-server=$(IMAGE):$(TAG) -n sumwhere

push:
	@echo $(DOCKER_PASSWORD) | docker login -u $(DOCKER_USERNAME) --password-stdin
	docker push $(IMAGE):$(VERSIONS)