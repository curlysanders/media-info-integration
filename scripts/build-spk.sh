#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
OUTPUT="$ROOT/dist"
STAGE="$OUTPUT/stage"
VERSION=$(sed -n 's/^version="\(.*\)"$/\1/p' "$ROOT/package/INFO")
ARCH=$(sed -n 's/^arch="\(.*\)"$/\1/p' "$ROOT/package/INFO")

rm -rf "$STAGE"
mkdir -p "$STAGE/package/bin" "$STAGE/package/nginx" "$STAGE/conf" "$STAGE/scripts"

docker run --rm \
  --user "$(id -u):$(id -g)" \
  -v "$ROOT:/src:ro" \
  -v "$STAGE:/out" \
  -w /src \
  golang:1.25.4-bookworm \
  sh -ceu 'export GOCACHE=/tmp/go-build GOMODCACHE=/tmp/go-mod; CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v1 go build -trimpath -buildvcs=false -ldflags="-s -w" -o /out/package/bin/runtime-probe ./cmd/runtime-probe; go run ./cmd/package-icon -small /out/PACKAGE_ICON.PNG -large /out/PACKAGE_ICON_256.PNG'

cp "$ROOT/package/INFO" "$STAGE/INFO"
cp "$ROOT/package/conf/privilege" "$STAGE/conf/privilege"
cp "$ROOT/package/conf/resource" "$STAGE/conf/resource"
cp -R "$ROOT/package/nginx/." "$STAGE/package/nginx/"
cp "$ROOT/package/scripts/start-stop-status" "$STAGE/scripts/start-stop-status"
cp "$ROOT/package/scripts/preupgrade" "$STAGE/scripts/preupgrade"
cp "$ROOT/package/scripts/postupgrade" "$STAGE/scripts/postupgrade"
chmod 0755 "$STAGE/scripts/start-stop-status" "$STAGE/scripts/preupgrade" "$STAGE/scripts/postupgrade" "$STAGE/package/bin/runtime-probe"
(cd "$STAGE/package" && tar -czf ../package.tgz bin nginx)
printf 'checksum="%s"\n' "$(md5sum "$STAGE/package.tgz" | awk '{print $1}')" >> "$STAGE/INFO"
(cd "$STAGE" && tar -cf "$OUTPUT/mediainfo-poc-$VERSION-$ARCH.spk" INFO PACKAGE_ICON.PNG PACKAGE_ICON_256.PNG conf scripts package.tgz)
sha256sum "$OUTPUT/mediainfo-poc-$VERSION-$ARCH.spk" > "$OUTPUT/mediainfo-poc-$VERSION-$ARCH.spk.sha256"
tar -tf "$OUTPUT/mediainfo-poc-$VERSION-$ARCH.spk"
