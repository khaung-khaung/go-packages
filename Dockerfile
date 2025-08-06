# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-packages-app

# Stage 2: Create the final lightweight image
FROM alpine:3.19
WORKDIR /app

# Override Alpine repositories to use HTTP
RUN echo "http://dl-cdn.alpinelinux.org/alpine/v3.19/main" > /etc/apk/repositories && \
    echo "http://dl-cdn.alpinelinux.org/alpine/v3.19/community" >> /etc/apk/repositories

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Copy the compiled binary
COPY --from=builder /go-packages-app /app/go-packages-app

CMD ["/app/go-packages-app"]