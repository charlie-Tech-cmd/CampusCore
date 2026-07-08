# =========================================================================
# Stage 1: Build the application
# =========================================================================
FROM golang:1.23-alpine AS builder

# Install packages needed during the build
RUN apk add --no-cache git ca-certificates build-base

# Set the working directory
WORKDIR /app

# Copy dependency files and download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application binary
# CGO_ENABLED=0 creates a static binary
# GOOS=linux builds for Linux
# -s -w removes debug information to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/bin/campuscore \
    cmd/server/main.go

# =========================================================================
# Stage 2: Create the runtime image
# =========================================================================
FROM alpine:3.19 AS runner

# Install certificates and timezone data
RUN apk add --no-cache ca-certificates tzdata

# Create a non-root user for security
RUN addgroup -S campusgroup && adduser -S campususer -G campusgroup

# Set the working directory
WORKDIR /home/campususer
USER campususer

# Copy the compiled binary from the build stage
COPY --from=builder --chown=campususer:campusgroup /app/bin/campuscore .

# The application listens on port 8080
EXPOSE 8080

# Start the application
ENTRYPOINT ["./campuscore"]