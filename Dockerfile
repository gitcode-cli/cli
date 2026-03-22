# Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build arguments for version info
ARG VERSION=dev
ARG COMMIT=unknown
ARG DATE=unknown

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" \
    -o /gc ./cmd/gc

# Final stage
FROM alpine:3.19

RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    git \
    bash \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1000 gc && \
    adduser -u 1000 -G gc -s /bin/sh -D gc

WORKDIR /home/gc

# Copy binary from builder
COPY --from=builder /gc /usr/local/bin/gc

# Copy completions
COPY --from=builder /app/completions /usr/share/completions

# Set ownership
RUN chown -R gc:gc /home/gc

# Switch to non-root user
USER gc

# Set environment
ENV PATH="/usr/local/bin:${PATH}"
ENV GC_PAGER=less

# Default command
ENTRYPOINT ["gc"]
CMD ["--help"]

# Labels
LABEL org.opencontainers.image.title="GitCode CLI"
LABEL org.opencontainers.image.description="Command line tool for GitCode"
LABEL org.opencontainers.image.url="https://gitcode.com"
LABEL org.opencontainers.image.source="https://github.com/gitcode-com/gitcode-cli"
LABEL org.opencontainers.image.vendor="GitCode"