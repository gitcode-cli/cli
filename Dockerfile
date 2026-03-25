# GitCode CLI Docker Image
# Uses pre-built binary from GoReleaser

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

# Copy pre-built binary (provided by GoReleaser)
COPY gc /usr/local/bin/gc

# Copy completions
COPY completions /usr/share/completions

# Set ownership
RUN chown -R gc:gc /home/gc && \
    chmod +x /usr/local/bin/gc

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