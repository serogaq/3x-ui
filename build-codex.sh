#!/usr/bin/env bash
set -euo pipefail

TARGETARCH=${TARGETARCH:-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')}
BUILD_WITH_ANTIZAPRET=${BUILD_WITH_ANTIZAPRET:-0}
WORK_DIR=${WORK_DIR:-$(pwd)}
export DEBIAN_FRONTEND=noninteractive

cd $WORK_DIR

export CGO_ENABLED=1
export GOOS=linux
export GOARCH=$TARGETARCH
if [ "$TARGETARCH" == "arm64" ]; then
  export GOARCH=arm64
  export CC=aarch64-linux-gnu-gcc
fi
go build -ldflags "-w -s" -o xui-build -v main.go

cp xui-build x-ui/
cp x-ui.service x-ui/
cp x-ui.sh x-ui/
mv x-ui/xui-build x-ui/x-ui
chmod +x x-ui/x-ui

cd x-ui
pkill x-ui && rm -f x-ui.log
touch x-ui.log && ./x-ui run > x-ui.log 2>&1 &
sleep 1 && cat x-ui.log && echo "Http Status Code: $(curl -s -o /dev/null -w %{http_code} http://localhost:2053)"
