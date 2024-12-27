# Use a multi-stage build
# Stage 1: Builder (for compiling, if needed)
FROM openjdk:11-jdk-slim AS builder

WORKDIR /app

# Copy only the necessary files (if any)
# COPY ./java-src /app/java-src

# Compile as a non-root user (if you need compilation as a separate step)
# RUN useradd -m appuser && chown -R appuser /app
# USER appuser
# RUN javac Main.java

# Stage 2: Runtime
FROM openjdk:11-jre-slim

# Create a non-root user
RUN groupadd -r appuser && useradd -r -g appuser appuser

WORKDIR /app

# Copy only the compiled class file(s) from the builder stage
# COPY --from=builder /app/Main.class /app/Main.class

# Set the non-root user
USER appuser

CMD ["bash"]