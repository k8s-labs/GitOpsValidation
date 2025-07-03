# syntax=docker/dockerfile:1.4
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Cache go mod downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gov ./cmd/gov

# -- Final image --
FROM alpine:3.20
WORKDIR /app

# Add non-root user
RUN adduser -D -u 10001 govuser

# Copy binary
COPY --from=builder /app/gov /app/gov

# Set permissions
RUN chown govuser:govuser /app/gov
USER govuser
ENV GOV_DAEMON true
EXPOSE 8080

ENTRYPOINT ["/app/gov"]
