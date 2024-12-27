# Use a multi-stage build
# Stage 1: Builder
FROM gcc:latest AS builder

# Create a non-root user
RUN groupadd -r appuser && useradd -r -g appuser appuser

WORKDIR /app

# Copy only the necessary files (if any)
# COPY ./cpp-src /app/cpp-src

# Compile as the non-root user (if needed for building from source)
# USER appuser
# RUN g++ -o myprogram myprogram.cpp

# Stage 2: Runtime
FROM gcr.io/distroless/cc-debian11

# Copy only the compiled binary from the builder stage
COPY --from=builder /usr/bin/g++ /usr/bin/g++

# Set a non-root user
USER 1000

WORKDIR /app

CMD ["bash"]