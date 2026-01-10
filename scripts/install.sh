#!/bin/bash
echo "Installing tunnerse CLI and Server..."

cd "$(dirname "$0")/.."

BIN_CLI="tunnerse"
BIN_SERVER="tunnerse-server"
BIN_DIR="bin"
SERVICE_FILE="tunnerse-server.service"

echo "Checking for compiled binaries..."
if [ ! -f "$BIN_DIR/$BIN_CLI" ]; then
    echo "ERROR: $BIN_CLI not found in $BIN_DIR/"
    echo "Please run ./scripts/build.sh first to compile the binaries."
    exit 1
fi

if [ ! -f "$BIN_DIR/$BIN_SERVER" ]; then
    echo "ERROR: $BIN_SERVER not found in $BIN_DIR/"
    echo "Please run ./scripts/build.sh first to compile the binaries."
    exit 1
fi

echo "Installing binaries..."
sudo mkdir -p /usr/local/bin

sudo cp "$BIN_DIR/$BIN_CLI" /usr/local/bin/
sudo cp "$BIN_DIR/$BIN_SERVER" /usr/local/bin/

sudo chmod +x /usr/local/bin/"$BIN_CLI"
sudo chmod +x /usr/local/bin/"$BIN_SERVER"

echo "Installing runtime assets into ~/.tunnerse/..."
TUNNERSE_HOME="$HOME/.tunnerse"
sudo mkdir -p "$TUNNERSE_HOME"

if [ -d "static" ]; then
    sudo rm -rf "$TUNNERSE_HOME/static" 2>/dev/null || true
    sudo cp -r "static" "$TUNNERSE_HOME/"
else
    echo "Warning: static/ directory not found at project root"
fi

sudo chown -R "$USER":"$USER" "$TUNNERSE_HOME" 2>/dev/null || true

echo "Configuring systemd service..."
# Cria o arquivo de serviço substituindo %i pelo usuário atual
sed "s/%i/$USER/g" scripts/"$SERVICE_FILE" > /tmp/"$SERVICE_FILE"

# Copia para o diretório de serviços do systemd
sudo cp /tmp/"$SERVICE_FILE" /etc/systemd/system/
sudo chmod 644 /etc/systemd/system/"$SERVICE_FILE"
rm /tmp/"$SERVICE_FILE"

# Recarrega o systemd para reconhecer o novo serviço
sudo systemctl daemon-reload

echo "Successfully installed both CLI and Server."
echo ""
echo "To manage the tunnerse-server daemon:"
echo "  Start:   sudo systemctl start tunnerse-server"
echo "  Stop:    sudo systemctl stop tunnerse-server"
echo "  Status:  sudo systemctl status tunnerse-server"
echo "  Enable:  sudo systemctl enable tunnerse-server  (start on boot)"
echo "  Logs:    sudo journalctl -u tunnerse-server -f"
echo ""
echo "Use 'tunnerse help' for CLI details."
echo
