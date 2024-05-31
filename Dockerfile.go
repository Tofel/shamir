# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.22 AS builder

# Set build arguments
ARG TARGETOS
ARG TARGETARCH

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the Go source code
COPY . .

# Build the Go application for the specified architecture
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o shamir shamir.go

# Use a minimal image to run the Go application for the specified platform
FROM --platform=$TARGETPLATFORM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go application from the builder stage
COPY --from=builder /app/shamir .

# Ensure the binary has execute permissions
RUN chmod +x shamir

# Define the entrypoint command
ENTRYPOINT ["./shamir"]