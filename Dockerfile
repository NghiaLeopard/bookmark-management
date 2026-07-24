FROM golang:1.25.0-alpine AS base

WORKDIR /opt/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

FROM base AS build-base

RUN go build -o bookmark-management ./cmd/api/main.go

FROM base AS test-exec

ARG _outputdir="/tmp/coverage"
ARG COVERAGE_EXCLUDE

RUN mkdir -p ${_outputdir} && \
    go test ./... -coverprofile=${_outputdir}/coverage.tmp -covermode=atomic -coverpkg=./... -p 1 && \
	grep -vE "${COVERAGE_EXCLUDE}" ${_outputdir}/coverage.tmp > ${_outputdir}/coverage.out && \
	go tool cover -html=${_outputdir}/coverage.out -o ${_outputdir}/coverage.html

FROM scratch as test

ARG _outputdir="/tmp/coverage"

COPY --from=test-exec ${_outputdir}/coverage.html /
COPY --from=test-exec ${_outputdir}/coverage.out /


FROM alpine:latest

WORKDIR /app

COPY --from=builder /opt/app/bookmark-management /app/bookmark-management
COPY --from=builder /opt/app/docs /app/docs


CMD ["./bookmark-management"]