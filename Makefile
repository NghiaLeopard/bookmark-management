.PHONY: run test dev-run swag docker-test docker-build dev-run-docker

GIT_TAG := $(shell git describe --tags --exact-match --abbrev=0 2>/dev/null)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
IMG_TAG := latest


IMG_NAME=ebvn/test_k04


ifneq ($(GIT_TAG),)
   IMG_TAG := $(GIT_TAG)
endif


export IMG_TAG


run:
	go run ./cmd/api/main.go

swag: 
	swag init -g ./cmd/api/main.go -o ./docs

COVERAGE_EXCLUDE = mocks|main.go|test|redis|docs
COVERAGE_THRESHOLD = 80

test:
	go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if awk -v t="$$total" -v th="$(COVERAGE_THRESHOLD)" 'BEGIN {exit !(t < th)}'; then \
		echo "Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

COVERAGE_FOLDER = ./test-output

docker-test:
	mkdir -p $(COVERAGE_FOLDER)
	docker buildx build --progress=plain --build-arg COVERAGE_EXCLUDE="$(COVERAGE_EXCLUDE)" --target test -t test:test --output ./test-output .
	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if awk -v t="$$total" -v th="$(COVERAGE_THRESHOLD)" 'BEGIN {exit !(t < th)}'; then \
		echo "Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

dev-run: swag run

docker-build:
	docker build -t $(IMG_NAME):$(IMG_TAG) .

DOCKER_USERNAME ?=
DOCKER_PASSWORD ?=

docker-login:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin

docker-release:
	docker push $(IMG_NAME):$(IMG_TAG)