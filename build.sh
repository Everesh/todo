#!/bin/bash
# Build script for Unix-like systems
# Usage: ./build.sh [all|current]

set -e

APP_NAME="todo"
BUILD_DIR="build"
VERSION="${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')}"
GO_FLAGS="-ldflags \"-X main.version=${VERSION}\""

build_platform() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT_NAME=$3
    
    echo "Building ${APP_NAME} for ${GOOS}/${GOARCH}..."
    GOOS=${GOOS} GOARCH=${GOARCH} go build ${GO_FLAGS} -o "${BUILD_DIR}/${OUTPUT_NAME}" ./cmd/todo
}

# Create build directory
mkdir -p "${BUILD_DIR}"

if [ "${1:-current}" = "all" ]; then
    echo "Building ${APP_NAME} for all platforms..."
    
    build_platform "windows" "amd64" "${APP_NAME}-windows-amd64.exe"
    build_platform "windows" "386" "${APP_NAME}-windows-386.exe"
    build_platform "windows" "arm64" "${APP_NAME}-windows-arm64.exe"
    build_platform "darwin" "amd64" "${APP_NAME}-darwin-amd64"
    build_platform "darwin" "arm64" "${APP_NAME}-darwin-arm64"
    build_platform "linux" "amd64" "${APP_NAME}-linux-amd64"
    build_platform "linux" "386" "${APP_NAME}-linux-386"
    build_platform "linux" "arm64" "${APP_NAME}-linux-arm64"
    build_platform "linux" "arm" "${APP_NAME}-linux-arm"
    
    echo ""
    echo "Done! Binaries are in ${BUILD_DIR}/"
else
    echo "Building ${APP_NAME} for current platform..."
    CURRENT_OS=$(go env GOOS)
    CURRENT_ARCH=$(go env GOARCH)
    EXTENSION=""
    [ "${CURRENT_OS}" = "windows" ] && EXTENSION=".exe"
    
    build_platform "${CURRENT_OS}" "${CURRENT_ARCH}" "${APP_NAME}${EXTENSION}"
    echo ""
    echo "Done! Binary is in ${BUILD_DIR}/"
fi

