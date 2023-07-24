.PHONY: image force-image build

bin := twelvedata-exporter
src := $(wildcard *.go)

# Default target
${bin}: Makefile ${src}
	go build -v -o "${bin}"

# Docker targets
image:
	docker build -t ${USER}/twelvedata-exporter .

force-image:
	docker build --no-cache -t ${USER}/twelvedata-exporter .
