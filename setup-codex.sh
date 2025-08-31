#!/usr/bin/env bash
set -euo pipefail

TARGETARCH=${TARGETARCH:-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')}
BUILD_WITH_ANTIZAPRET=${BUILD_WITH_ANTIZAPRET:-0}
GO_VERSION=${GO_VERSION:-1.25.0}
export DEBIAN_FRONTEND=noninteractive
export CGO_ENABLED=1
export GOOS=linux

export_postman_collection() {
  local id="$1"
  local out="$2"
  local url="https://www.postman.com/collections/${id}?format=2.1"
  curl -sSL "$url" -o "$out"
}

apt update
apt install --assume-yes curl wget unzip procps file bsdmainutils busybox libc6-dev

go install "golang.org/dl/go${GO_VERSION}@latest"
"go${GO_VERSION}" download

go version

go mod download

TOOLCHAIN_URL=""
MUSL_CC_HOST="https://github.com/musl-cc/musl.cc/releases/download/v0.0.1"
if [ "$TARGETARCH" == "amd64" ]; then
  TOOLCHAIN_URL="$MUSL_CC_HOST/x86_64-linux-musl-native.tgz"
elif [ "$TARGETARCH" == "arm64" ]; then
  TOOLCHAIN_URL="$MUSL_CC_HOST/aarch64-linux-musl-native.tgz"
fi
echo "Downloading musl toolchain for ${TARGETARCH}"
curl -LO "$TOOLCHAIN_URL"
tar -xf *.tgz
TOOLCHAIN_DIR=$(find . -maxdepth 1 -type d -name "*-native" | head -n1)
TOOLCHAIN_DIR=$(realpath "$TOOLCHAIN_DIR")
export PATH="$TOOLCHAIN_DIR/bin:$PATH"
export CC=$(find $TOOLCHAIN_DIR/bin -name '*-gcc' | head -n1)
chmod +x $CC
echo "Using CC=$CC"

mkdir -p x-ui/bin
cd x-ui/bin

Xray_URL="https://github.com/XTLS/Xray-core/releases/download/v25.8.29/"
if [ "$TARGETARCH" == "amd64" ]; then
  wget -q ${Xray_URL}Xray-linux-64.zip
  unzip Xray-linux-64.zip
  rm -f Xray-linux-64.zip
elif [ "$TARGETARCH" == "arm64" ]; then
  wget -q ${Xray_URL}Xray-linux-arm64-v8a.zip
  unzip Xray-linux-arm64-v8a.zip
  rm -f Xray-linux-arm64-v8a.zip
fi
rm -f geoip.dat geosite.dat
wget -q https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat
wget -q https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat
wget -q -O geoip_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geoip.dat
wget -q -O geosite_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geosite.dat
wget -q -O geoip_RU.dat https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geoip.dat
wget -q -O geosite_RU.dat https://github.com/runetfreedom/russia-v2ray-rules-dat/releases/latest/download/geosite.dat
mv xray xray-linux-$TARGETARCH

cd ../..

export_postman_collection "5146551-dda3cab3-0e33-485f-96f9-d4262f437ac5" "postman_api_collection.json"
