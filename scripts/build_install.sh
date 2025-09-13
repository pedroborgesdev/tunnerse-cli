#!/bin/bash
echo "Installing tunnerse..."

cd ..
cd src

BIN_NAME="tunnerse"

echo "Compiling with go build..."
go build -o "$BIN_NAME"

if [ $? -ne 0 ]; then
    echo "Compilation failed."
    exit 1
fi

mkdir -p /usr/local/bin

mv "$BIN_NAME" /usr/local/bin/

chmod +x /usr/local/bin/"$BIN_NAME"

echo "Successfully installed. Use 'tunnerse help' for more details."
echo
