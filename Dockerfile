# Dockerfile for twelvedata-exporter

FROM scratch

# dockers_v2 lays the build context out as linux/<arch>/<binary>
ARG TARGETPLATFORM

# Copy ca-certificates for HTTPS requests to the Twelve Data API
COPY --from=alpine:3.24.2@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the pre-built binary from GoReleaser
COPY $TARGETPLATFORM/twelvedata-exporter /twelvedata-exporter

# extra_files in .goreleaser.yml is what puts these in the build context
COPY LICENSE NOTICE /

# Create a non-root user (using numeric ID for scratch image)
USER 65534:65534

# Declare the port; publishing it still requires docker run -p
EXPOSE 10016

# Set the entrypoint. No CMD: a default argument of --help would make the
# documented `docker run … ghcr.io/umatare5/twelvedata-exporter` print help and
# exit instead of starting the exporter.
ENTRYPOINT ["/twelvedata-exporter"]
