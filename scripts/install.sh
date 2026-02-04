#!/bin/bash

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "Error: This script must be run with sudo"
    echo "Usage: sudo ./install.sh"
    exit 1
fi

echo "Installing tunnerse CLI and Server..."

SCRIPT_DIR="$(dirname "$0")"
cd "$SCRIPT_DIR"

BIN_CLI="tunnerse"
BIN_SERVER="tunnerse-server"

# Verifica se os binários existem
if [ ! -f "$BIN_CLI" ]; then
    echo "Error: $BIN_CLI not found. Please compile first with: go build -o $BIN_CLI ../cmd/cli"
    exit 1
fi

if [ ! -f "$BIN_SERVER" ]; then
    echo "Error: $BIN_SERVER not found. Please compile first with: go build -o $BIN_SERVER ../cmd/server"
    exit 1
fi

echo "Installing binaries to /usr/local/bin/..."
mkdir -p /usr/local/bin

cp "$BIN_CLI" /usr/local/bin/
cp "$BIN_SERVER" /usr/local/bin/

chmod +x /usr/local/bin/"$BIN_CLI"
chmod +x /usr/local/bin/"$BIN_SERVER"

# Install systemd service
echo "Installing systemd service..."

# Pega o usuário real (mesmo quando rodado com sudo)
REAL_USER="${SUDO_USER:-$USER}"

# Cria o arquivo de serviço com o usuário correto
cat > /tmp/tunnerse-server.service << EOF
[Unit]
Description=Tunnerse Server - Local tunnel management daemon
After=network.target

[Service]
Type=simple
User=$REAL_USER
Group=$REAL_USER
ExecStart=/usr/local/bin/tunnerse-server
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

# Move o serviço para o systemd
mv /tmp/tunnerse-server.service /etc/systemd/system/
systemctl daemon-reload

echo ""
echo "Systemd service installed successfully for user: $REAL_USER"
echo "Service will store data in: /home/$REAL_USER/.tunnerse/"
echo ""
echo "To manage the service:"
echo "  sudo systemctl enable tunnerse-server    # Enable on boot"
echo "  sudo systemctl start tunnerse-server     # Start now"
echo "  sudo systemctl stop tunnerse-server      # Stop"
echo "  sudo systemctl status tunnerse-server    # Check status"
echo "  sudo journalctl -u tunnerse-server -f    # View logs"

echo ""
echo "Successfully installed both CLI and Server."
echo "Use 'tunnerse help' for CLI details."
echo "Use 'tunnerse-server' to start the server manually."
echo "Or use systemd to run as a service (see above)."
echo
