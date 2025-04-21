# Stage 1: Build
FROM --platform=linux/amd64 golang:1.23 AS builder
WORKDIR /app
COPY . .
# Build the Go binary for Linux AMD64 platform
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/say-more-server .

# Stage 2: Runtime
FROM --platform=linux/amd64 alpine:latest
# Create a working directory
RUN mkdir -p /app
WORKDIR /app
# Copy the compiled binary and set executable permissions
COPY --from=builder --chmod=755 /app/say-more-server /app/
# Run the binary directly (avoid using sh -c)
CMD ["/app/say-more-server", "start"]