FROM golang:1.18-alpine AS build_base

RUN apk add --no-cache git build-base

# Set the Current Working Directory inside the container
WORKDIR /build

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/connector-api .

# Start fresh from a smaller image
FROM alpine:3.9
RUN apk add ca-certificates

COPY --from=build_base /build/out/connector-api /app/connector-api

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["/app/connector-api"]