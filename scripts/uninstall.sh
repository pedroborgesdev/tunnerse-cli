#!/bin/bash
echo "Uninstalling tunnerse CLI and Server..."

SERVICE_FILE="tunnerse-server.service"

# Para o serviço se estiver rodando
if systemctl is-active --quiet tunnerse-server; then
    echo "Stopping tunnerse-server service..."
    sudo systemctl stop tunnerse-server
fi

# Desabilita o serviço se estiver habilitado
if systemctl is-enabled --quiet tunnerse-server 2>/dev/null; then
    echo "Disabling tunnerse-server service..."
    sudo systemctl disable tunnerse-server
fi

# Remove o arquivo de serviço
if [ -f "/etc/systemd/system/$SERVICE_FILE" ]; then
    echo "Removing systemd service file..."
    sudo rm /etc/systemd/system/"$SERVICE_FILE"
    sudo systemctl daemon-reload
fi

# Remove os binários
echo "Removing binaries..."
sudo rm -f /usr/local/bin/tunnerse
sudo rm -f /usr/local/bin/tunnerse-server

# Remove os dados (opcional - comentado por segurança)
# echo "Removing data files..."
# rm -rf ~/.tunnerse

echo "Uninstall complete."
echo ""
echo "Note: Tunnel logs and database were NOT removed."
echo "To remove them manually, delete: /usr/local/bin/tunnels/"
