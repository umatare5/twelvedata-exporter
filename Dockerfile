# Dockerfile for twelvedata-exporter

FROM scratch

# Copy ca-certificates for HTTPS requests to twelvedata-exporter controllers
COPY --from=alpine:latest@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the pre-built binary from GoReleaser
COPY twelvedata-exporter /twelvedata-exporter

# Create a non-root user (using numeric ID for scratch image)
USER 65534:65534

# Set the entrypoint
ENTRYPOINT ["/twelvedata-exporter"]

# Default command shows help
CMD ["--help"]
