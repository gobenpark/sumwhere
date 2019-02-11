GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
IMAGE=qkrqjadn/sumwhere
DOCKER_PASSWORD=1q2w3e4r
DOCKER_USERNAME=qkrqjadn
VERSION=1.0
GITCOMMITCOUNT:=$$(git rev-list HEAD | wc -l | tr -d ' ')
GITHASH:=$$(git rev-parse --short HEAD)
DATETIME:=$$(date "+%Y%m%d")
VERSIONS:=$(VERSION).$(GITCOMMITCOUNT)-$(GITHASH)-$(DATETIME)
#https://codecov.io/
.PHONY: clean docker-build rolling-update sumwhere test


clean:
	$(GOCLEAN)

sumwhere:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o $@ -ldflags "-X main.ServiceVersion=$(VERSIONS)" *.go

docker-build:
	@docker build -t $(IMAGE):$(VERSIONS) .

rolling-update:
	ssh -o StrictHostKeyChecking=no root@210.100.177.146 -p 55555 kubectl set image deployment/sumwhere-server sumwhere-server=$(IMAGE):$(VERSIONS) -n sumwhere

push:
	@echo $(DOCKER_PASSWORD) | docker login -u $(DOCKER_USERNAME) --password-stdin
	docker push $(IMAGE):$(VERSIONS)