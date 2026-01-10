#!/bin/bash

# Script para gerenciar o serviço tunnerse-server

USER_NAME=$(whoami)

case "$1" in
    enable)
        echo "Enabling tunnerse-server service for user $USER_NAME..."
        sudo systemctl enable tunnerse-server@$USER_NAME
        echo "Service enabled. It will start automatically on boot."
        ;;
    disable)
        echo "Disabling tunnerse-server service for user $USER_NAME..."
        sudo systemctl disable tunnerse-server@$USER_NAME
        echo "Service disabled."
        ;;
    start)
        echo "Starting tunnerse-server service for user $USER_NAME..."
        sudo systemctl start tunnerse-server@$USER_NAME
        echo "Service started."
        ;;
    stop)
        echo "Stopping tunnerse-server service for user $USER_NAME..."
        sudo systemctl stop tunnerse-server@$USER_NAME
        echo "Service stopped."
        ;;
    restart)
        echo "Restarting tunnerse-server service for user $USER_NAME..."
        sudo systemctl restart tunnerse-server@$USER_NAME
        echo "Service restarted."
        ;;
    status)
        sudo systemctl status tunnerse-server@$USER_NAME
        ;;
    logs)
        sudo journalctl -u tunnerse-server@$USER_NAME -f
        ;;
    *)
        echo "Usage: $0 {enable|disable|start|stop|restart|status|logs}"
        echo ""
        echo "Commands:"
        echo "  enable   - Enable service to start on boot"
        echo "  disable  - Disable service from starting on boot"
        echo "  start    - Start the service"
        echo "  stop     - Stop the service"
        echo "  restart  - Restart the service"
        echo "  status   - Show service status"
        echo "  logs     - Show service logs (real-time)"
        echo ""
        echo "Current user: $USER_NAME"
        echo "Service name: tunnerse-server@$USER_NAME"
        exit 1
        ;;
esac
