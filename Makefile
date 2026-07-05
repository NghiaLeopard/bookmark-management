.PHONY: run 

run:
	go run ./cmd/api/main.go


# -coverpkg: cover the packages that are used in the tests
# -covermode: atomic, count, set, or default
# -p 1: run the tests in parallel
# -coverprofile: generate the coverage profile
# -coverhtml: generate the coverage HTML report

COVERAGE_EXCLUDE=mocks|main
test:
	go test ./... -coverprofile=coverage.tmp -coverpkg=./... -covermode=atomic -p 1
	grep -vE "${COVERAGE_EXCLUDE}" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html 

