.PHONY: run test dev-run

run:
	go run ./cmd/api/main.go

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

dev-run: run test
