# Multi-stage Dockerfile for Skyclerk
# Builds both frontend Angular apps and Go backend into a single image

# Stage 1: Build Frontend
FROM node:16-alpine AS frontend-builder

# Set working directory
WORKDIR /app

# Build main frontend app
COPY frontend/package*.json ./frontend/
RUN cd frontend && npm ci

COPY frontend/ ./frontend/
RUN cd frontend && npm run build -- --configuration=production

# Build centcom admin app  
COPY centcom/package*.json ./centcom/
RUN cd centcom && npm ci --legacy-peer-deps

COPY centcom/ ./centcom/
RUN cd centcom/src && export NODE_ENV=production && npx tailwindcss build tailwind.css -o styles.css
RUN cd centcom && npm run build -- --configuration=production --base-href="/centcom/"

# Stage 2: Build Go Backend
FROM golang:1.21-bullseye AS backend-builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    libc6-dev \
    libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source
COPY backend/ ./

# Build the application
RUN CGO_ENABLED=1 go build -o skyclerk .

# Stage 3: Extract imaginary binary from Docker image
FROM --platform=linux/amd64 h2non/imaginary:1.1.0 AS imaginary-extractor

# Stage 4: Runtime Image with Imaginary
FROM debian:bullseye-slim

# Install runtime dependencies including supervisor for multi-process management and libvips
RUN apt-get update && apt-get install -y \
    ca-certificates \
    sqlite3 \
    supervisor \
    libvips42 \
    && rm -rf /var/lib/apt/lists/*

# Copy imaginary binary from the official Docker image (it's at /bin/imaginary in the image)
COPY --from=imaginary-extractor /bin/imaginary /usr/local/bin/imaginary

# Create app directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=backend-builder /app/skyclerk .

# Copy frontend build files
COPY --from=frontend-builder /app/frontend/dist/frontend ./frontend
COPY --from=frontend-builder /app/centcom/dist ./centcom

# Copy fonts directory
COPY fonts/ ./fonts/

# Create directory for SQLite database
RUN mkdir -p /app/data

# Create supervisor configuration
RUN echo '[supervisord]' > /etc/supervisor/conf.d/skyclerk.conf && \
    echo 'nodaemon=true' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo 'user=root' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo '' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo '[program:skyclerk]' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo 'command=/app/skyclerk' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo 'directory=/app' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo 'autostart=true' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo 'autorestart=true' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo 'stdout_logfile=/dev/stdout' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo 'stdout_logfile_maxbytes=0' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo 'stderr_logfile=/dev/stderr' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo 'stderr_logfile_maxbytes=0' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo '' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo '[program:imaginary]' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo 'command=/usr/local/bin/imaginary -concurrency 50 -enable-url-source -p 9000' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo 'autostart=true' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo 'autorestart=true' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo 'stdout_logfile=/dev/stdout' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo 'stdout_logfile_maxbytes=0' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo 'stderr_logfile=/dev/stderr' >> /etc/supervisor/conf.d/skyclerk.conf && \
    echo 'stderr_logfile_maxbytes=0' >> /etc/supervisor/conf.d/skyclerk.conf

# Expose ports for both services
EXPOSE 8080

# Set environment variables
ENV HTTP_PORT=8080
ENV IMAGINARY_HOST=http://127.0.0.1:9000

# We set to local because we are not doing https with this app.
ENV APP_ENV=local

# Run both services via supervisor
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/skyclerk.conf"]