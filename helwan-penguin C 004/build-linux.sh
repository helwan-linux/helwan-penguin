```bash
#!/usr/bin/env bash
set -e

cd "$(dirname "$0")"

echo "========================================"
echo "       Helwan Penguin Linux Build"
echo "========================================"

echo "[1/4] Checking Go..."
go version

echo "[2/4] Formatting Go files..."
find cmd internal -type f -name '*.go' -exec gofmt -w {} +

echo "[3/4] Running tests..."
go test ./...

echo "[4/4] Building Linux executable..."
mkdir -p dist

GOOS=linux GOARCH=amd64 CGO_ENABLED=1 \
go build -trimpath -o dist/HelwanPenguin ./cmd/helwan-penguin

echo
echo "BUILD SUCCESSFUL"
echo "Executable: $(pwd)/dist/HelwanPenguin"
```
