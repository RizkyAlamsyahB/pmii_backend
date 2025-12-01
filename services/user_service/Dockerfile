# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

# Copy everything needed for build
COPY shared ./shared
COPY services/user_service ./services/user_service

# Build from service directory
WORKDIR /build/services/user_service

# Update replace path for Docker
RUN sed -i 's|../../shared|/build/shared|g' go.mod

# Download dependencies
RUN go mod download

# Build the application  
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/api

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy migration files
COPY services/user_service/migrations ./migrations

EXPOSE 8081

CMD ["./main"]
