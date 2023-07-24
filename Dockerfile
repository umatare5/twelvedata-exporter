FROM golang:1-alpine AS builder

ARG UID=60000
ARG TWELVEDATA_API_KEY
ENV TWELVEDATA_API_KEY=$TWELVEDATA_API_KEY

# Copy the repo contents into /tmp/build
WORKDIR /tmp/build
COPY . .

RUN cd /tmp/build && \
    go mod download && \
    go build

# Build the small image
FROM alpine
WORKDIR /app
COPY --from=builder /tmp/build/twelvedata-exporter .

EXPOSE 9341
USER ${UID}
ENTRYPOINT [ "./twelvedata-exporter" ]
