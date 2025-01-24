FROM golang:1.23-alpine as builder

ARG CGO_ENABLED=0
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy go.mod and go.sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
COPY . .

# Build the Go app
RUN go build -o frontendApp /app/cmd/frontend

FROM scratch

COPY --from=builder /app/frontendApp /frontendApp

ENTRYPOINT ["/frontendApp"]
