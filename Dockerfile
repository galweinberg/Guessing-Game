# ---------- Stage 1: Builder ----------
    FROM golang:1.20-alpine AS builder

    # Set environment variables
    ENV GO111MODULE=on \
        CGO_ENABLED=0 \
        GOOS=linux \
        GOARCH=amd64
    
    # Create app directory and copy files
    WORKDIR /app
    COPY . .
    
    # Build the Go binary
    RUN go build -o server .
    
    # ---------- Stage 2: Minimal runtime ----------
    FROM alpine:latest
    
    # Set working directory in the minimal image
    WORKDIR /app
    
    # Copy the compiled binary from builder
    COPY --from=builder /app/server .
    
    # Expose the port your server listens on
    EXPOSE 8080
    
    # Command to run the server
    CMD ["./server", "server"]
    