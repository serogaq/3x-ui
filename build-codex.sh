#!/usr/bin/env bash
set -euo pipefail

TARGETARCH=${TARGETARCH:-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')}
BUILD_WITH_ANTIZAPRET=${BUILD_WITH_ANTIZAPRET:-0}
WORK_DIR=${WORK_DIR:-$(pwd)}
GO_VERSION=${GO_VERSION:-1.24.5}
TOOLCHAIN_DIR=$(find . -maxdepth 2 -type d -name "*-cross" | head -n1)
TOOLCHAIN_DIR=$(realpath "$TOOLCHAIN_DIR")
export DEBIAN_FRONTEND=noninteractive
export CGO_ENABLED=1
export GOOS=linux
export GOARCH=$TARGETARCH
export PATH="$TOOLCHAIN_DIR/bin:$PATH"
export CC=$(find $TOOLCHAIN_DIR/bin -name '*-gcc' | head -n1)

cd $WORK_DIR

test -f config/custom_version || touch config/custom_version

ls -alR
#if [ "$TARGETARCH" == "arm64" ]; then
#  export GOARCH=arm64
#  export CC=aarch64-linux-gnu-gcc
#fi
#go build -ldflags "-w -s" -o xui-build -v main.go

echo "Using CC=$CC"

go build -ldflags "-w -s -linkmode external -extldflags '-static'" -o xui-build -v main.go

mv xui-build x-ui/
cp x-ui.service x-ui/
cp x-ui.sh x-ui/
mv x-ui/xui-build x-ui/x-ui
chmod +x x-ui/x-ui

cd x-ui
file x-ui
ldd x-ui || echo "Static binary confirmed"
pkill x-ui && rm -f x-ui.log
touch x-ui.log && ./x-ui run > x-ui.log 2>&1 &
sleep 3 && cat x-ui.log && echo "Http Status Code: $(curl -s -o /dev/null -w %{http_code} http://localhost:2053)"
