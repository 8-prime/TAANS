# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY /cmd ./cmd
COPY /internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o taans ./cmd/main.go

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/taans /taans

USER nonroot:nonroot

ENV PORT=8080

ENTRYPOINT ["/taans"]