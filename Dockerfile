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

# Stage 3: Runtime Image
FROM debian:bullseye-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

# Create app directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=backend-builder /app/skyclerk .

# Copy frontend build files
COPY --from=frontend-builder /app/frontend/dist/frontend ./frontend
COPY --from=frontend-builder /app/centcom/dist ./centcom

# Create directory for SQLite database
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release
ENV DB_PATH=/app/data/skyclerk.db
ENV HTTP_PORT=8080
ENV SITE_URL=http://localhost:8080
ENV HTTP_LOG_REQUESTS=true

# We set to local because we are not doing https with this app.
ENV APP_ENV=local

# Run the application
CMD ["./skyclerk"]