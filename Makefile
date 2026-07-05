.PHONY: run test

run:
	go run ./cmd/api/main.go

test:
	go test $$(go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./...) \
		-coverprofile=coverage.out -coverpkg=./... -covermode=atomic -p 1 -count=1
	go tool cover -html=coverage.out -o coverage.html
