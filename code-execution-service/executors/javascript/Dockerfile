FROM node:lts-slim

# Create a non-root user
RUN groupadd -r appuser && useradd -r -g appuser appuser

WORKDIR /app

# Set the non-root user
USER appuser

CMD ["bash"]