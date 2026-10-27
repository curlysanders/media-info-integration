#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
docker run --rm \
  --user "$(id -u):$(id -g)" \
  -v "$ROOT:/src" \
  -w /src \
  golang:1.25.4-bookworm \
  sh -ceu 'GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod go test ./...'
