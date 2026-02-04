#!/bin/bash
set -e

cd "$(dirname "$0")/.."

echo "Building tunnerse CLI and Server for Linux and Windows..."
echo ""

mkdir -p bin

# Build for Linux
echo "Building for Linux..."
echo "  - CLI..."
GOOS=linux GOARCH=amd64 go build -o bin/tunnerse ./cmd/cli
echo "  ✓ CLI built: bin/tunnerse"

echo "  - Server..."
GOOS=linux GOARCH=amd64 go build -o bin/tunnerse-server ./cmd/server
echo "  ✓ Server built: bin/tunnerse-server"

echo ""

# Build for Windows
echo "Building for Windows..."
echo "  - CLI..."
GOOS=windows GOARCH=amd64 go build -o bin/tunnerse.exe ./cmd/cli
echo "  ✓ CLI built: bin/tunnerse.exe"

echo "  - Server..."
GOOS=windows GOARCH=amd64 go build -o bin/tunnerse-server.exe ./cmd/server
echo "  ✓ Server built: bin/tunnerse-server.exe"

echo ""
echo "Build complete! Binaries are in the bin/ directory:"
echo ""
echo "Linux:"
echo "  - ./bin/tunnerse"
echo "  - ./bin/tunnerse-server"
echo ""
echo "Windows:"
echo "  - bin/tunnerse.exe"
echo "  - bin/tunnerse-server.exe"
echo ""
