FROM golang:1.23-alpine

# Install required packages
RUN apk add --no-cache tzdata ca-certificates

# Set timezone
ENV TZ=Asia/Dhaka \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked Go binary
RUN go build -o bin/main ./cmd/main.go

# Expose application port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["./bin/main"]
