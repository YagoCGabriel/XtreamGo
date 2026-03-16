#!/bin/bash
VERSION="v1.0.0"
LDFLAGS="-s -w -X main.Version=$VERSION"

mkdir -p dist

echo "Building $VERSION..."

GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -o dist/xtreamgo-windows-amd64.exe .
GOOS=linux   GOARCH=amd64 go build -ldflags="$LDFLAGS" -o dist/xtreamgo-linux-amd64 .
GOOS=darwin  GOARCH=amd64 go build -ldflags="$LDFLAGS" -o dist/xtreamgo-darwin-amd64 .
GOOS=darwin  GOARCH=arm64 go build -ldflags="$LDFLAGS" -o dist/xtreamgo-darwin-arm64 .

echo "Done. Files in ./dist:"
ls -lh dist/
