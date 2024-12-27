# Multi-stage build for Go
# Stage 1: Builder
FROM golang:1.20 AS builder

WORKDIR /app

# Copy only necessary files for building (if you are building here)
# COPY go.mod go.sum ./
# COPY *.go ./

# Build as a non-root user (if building from source)
# RUN useradd -m appuser && chown -R appuser /app
# USER appuser
# RUN go build -o my-go-app .

# Stage 2: Runtime
FROM gcr.io/distroless/static-debian11

# Create a non-root user
USER 1000

WORKDIR /app

# Copy the compiled Go binary from the builder stage
COPY --from=builder /usr/local/go/bin/go /usr/local/go/bin/go

CMD ["bash"]