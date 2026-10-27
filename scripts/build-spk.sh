#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
OUTPUT="$ROOT/dist"
STAGE="$OUTPUT/stage"
VERSION=$(sed -n 's/^version="\(.*\)"$/\1/p' "$ROOT/package/INFO")
ARCH=$(sed -n 's/^arch="\(.*\)"$/\1/p' "$ROOT/package/INFO")

rm -rf "$STAGE"
mkdir -p "$STAGE/package/bin" "$STAGE/package/nginx" "$STAGE/package/conf" "$STAGE/scripts"

docker run --rm \
  --user "$(id -u):$(id -g)" \
  -v "$ROOT:/src:ro" \
  -v "$STAGE/package/bin:/out" \
  -w /src \
  golang:1.25.4-bookworm \
  sh -ceu 'GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v1 go build -trimpath -buildvcs=false -ldflags="-s -w" -o /out/runtime-probe ./cmd/runtime-probe'

cp "$ROOT/package/INFO" "$STAGE/INFO"
cp -R "$ROOT/package/conf" "$STAGE/conf"
cp -R "$ROOT/package/nginx/." "$STAGE/package/nginx/"
cp "$ROOT/package/scripts/start-stop-status" "$STAGE/scripts/start-stop-status"
chmod 0755 "$STAGE/scripts/start-stop-status" "$STAGE/package/bin/runtime-probe"
(cd "$STAGE/package" && tar -czf ../package.tgz bin nginx)
(cd "$STAGE" && tar -czf "$OUTPUT/mediainfo-poc-$VERSION-$ARCH.spk" INFO conf scripts package.tgz)
sha256sum "$OUTPUT/mediainfo-poc-$VERSION-$ARCH.spk" > "$OUTPUT/mediainfo-poc-$VERSION-$ARCH.spk.sha256"
tar -tzf "$OUTPUT/mediainfo-poc-$VERSION-$ARCH.spk"
