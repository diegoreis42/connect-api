# Use an official Go image as the base builder
FROM golang:1.23-bookworm AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to leverage caching
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with optimizations
RUN go build -o /connect-api

# Use a minimal runtime image for the final container
FROM debian:bookworm-slim AS runtime

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /connect-api /app/connect-api

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/connect-api"]
