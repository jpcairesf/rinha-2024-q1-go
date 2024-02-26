# syntax=docker/dockerfile:1
##
## Build the application from source
##

FROM golang:1.22.0-alpine3.18 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o main main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rinha-2024-q1-go .

##
## Run the tests in the container
##

FROM build-stage AS run-test-stage
RUN go test -v ./...

##
## Deploy the application binary into a lean image
##

FROM amd64/alpine:3.18 AS build-release-stage

WORKDIR /app

COPY --from=build-stage /app/main .

EXPOSE 8080
CMD ["/app/main"]