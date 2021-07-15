VERSION = $$(cat VERSION)
IMG_NAME = gpaggi/ethapi

.PHONY: build build-docker clean distribute clean-image fmt load-test

build:
	@cd src; CGO_ENABLED=0 go build -v -ldflags "-X github.com/gpaggi/ethapi/version.Version=$(VERSION) -s -w" -o ../bin/ethapi .

clean:
	@rm ethapi

build-docker:
	@cd src; docker build --build-arg VERSION=$(VERSION) -t $(IMG_NAME):$(VERSION) -f ../Dockerfile .
	@cd src; docker tag $(IMG_NAME):$(VERSION) $(IMG_NAME):latest

distribute: build-docker
	@docker push $(IMG_NAME):$(VERSION)
	@docker push $(IMG_NAME):latest

clean-image:
	@docker image rm $(IMG_NAME):$(VERSION)

fmt:
	@cd src; go mod tidy
	@cd src; find . -type f -name '*.go' | xargs gofmt -w -s

load-test:
	@vegeta attack -duration=60s -rate=0 -max-workers=500 -targets=test/target.list | vegeta report