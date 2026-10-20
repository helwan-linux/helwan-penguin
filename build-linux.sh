#!/bin/bash

# الانتقال لمسار السكربت الحالي
cd "$(dirname "$0")"

echo
echo "========================================"
echo "        Helwan Penguin Build System     "
echo "========================================"
echo

echo "[1/4] Checking Go..."
go version
if [ $? -ne 0 ]; then
    echo "BUILD FAILED: Go is not installed."
    exit 1
fi

echo
echo "[2/4] Formatting source files..."
gofmt -w cmd internal
if [ $? -ne 0 ]; then
    echo "BUILD FAILED: gofmt encountered an error."
    exit 1
fi

echo
echo "[3/4] Running tests..."
go test ./...
if [ $? -ne 0 ]; then
    echo "BUILD FAILED: Tests failed."
    exit 1
fi

echo
echo "[4/4] Building Linux executable..."

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=1

go build -trimpath -o "packaging/arch/HelwanPenguin" ./cmd/helwan-penguin
if [ $? -ne 0 ]; then
    echo "BUILD FAILED: Compilation error."
    exit 1
fi

echo
echo "========================================"
echo "BUILD SUCCESSFUL"
echo "Executable:"
echo "$(pwd)/packaging/arch/HelwanPenguin"
echo "========================================"
echo