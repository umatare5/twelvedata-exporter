# Dockerfile for twelvedata-exporter

FROM scratch

# Copy ca-certificates for HTTPS requests to twelvedata-exporter controllers
COPY --from=alpine:latest@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the pre-built binary from GoReleaser
COPY twelvedata-exporter /twelvedata-exporter

# Create a non-root user (using numeric ID for scratch image)
USER 65534:65534

# Set the entrypoint
ENTRYPOINT ["/twelvedata-exporter"]

# Default command shows help
CMD ["--help"]
