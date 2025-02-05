# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Install ffmpeg
RUN apk add --no-cache ffmpeg

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and install the Go dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o worker cmd/worker/main.go

# Stage 2: Create non-root user in a complete image
FROM alpine:3.19.1 as security_provider

# Create a non-root group and user
RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot

# Install ffmpeg
RUN apk add --no-cache ffmpeg

# Stage 3: Run in a scratch image
FROM scratch as production

# Copy the /etc/passwd file from the previous stage
COPY --from=security_provider /etc/passwd /etc/passwd

# Set the non-root user
USER nonroot

# Copy only the binary from the build stage to the final image
COPY --from=builder /app/worker /

# Copy ffmpeg binary
COPY --from=security_provider /usr/bin/ffmpeg /usr/bin/ffmpeg

# Run the web usecases on container startup
CMD ["./worker"]