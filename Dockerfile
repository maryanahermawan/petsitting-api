# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY lib ./lib
COPY application.yaml ./application.yaml

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app-server ./cmd/app

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app-server /app-server
COPY --from=build /app/application.yaml /application.yaml

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app-server"]
