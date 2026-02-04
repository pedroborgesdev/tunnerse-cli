<div align="center">
	<img src="static/icon.png" alt="Tunnerse logo" width="96" height="96" />
	<h1>Tunnerse CLI</h1>
	<p><strong>Reverse tunnels for local development, demos, and testing.</strong></p>
</div>

> You can check whether the tunneling service is running at: https://tunnerse.com

---

## What is Tunnerse?

Tunnerse is a reverse tunneling system that exposes a local service to the public internet through a secure intermediary server. This repository contains the **Tunnerse CLI** and the **local tunnel server (tunnerse-server)** that run on your machine and connect to the Tunnerse API.

Use it to:
- Share local apps with teammates or clients
- Receive webhooks on localhost
- Demo prototypes without deploying
- Test integrations in real-world environments

## How it works

1. A public user hits your Tunnerse URL (e.g., `https://<tunnel>.tunnerse.com`).
2. The Tunnerse API queues the request for your tunnel.
3. Your local **tunnerse-server** fetches that request, forwards it to your local port, and captures the response.
4. The response is sent back through the API to the original client.

The CLI is a friendly wrapper around this flow, letting you create and manage tunnels with a few commands.

## Components

- **tunnerse** (CLI): creates and manages tunnels.
- **tunnerse-server** (local daemon): receives tunnel requests and forwards them to your local service.
- **Tunnerse API** (remote): public entry point that coordinates tunnel traffic.

## Quick start

### 1) Build from source

```bash
go build -o tunnerse ./cmd/cli
go build -o tunnerse-server ./cmd/server
```

### 2) Start the local daemon

```bash
./tunnerse-server
```

The local server listens on `http://localhost:9988` by default and stores data in `~/.tunnerse`.

### 3) Create a tunnel

**Quick (foreground):**

```bash
./tunnerse quick my-app 3000
```

**Persistent (background-managed):**

```bash
./tunnerse new my-app 3000
```

The CLI prints the public URL for your tunnel once it’s ready.

## CLI commands

| Command | Description |
| --- | --- |
| `tunnerse new <name> <port>` | Create a persistent tunnel (runs in background) |
| `tunnerse quick <name> <port>` | Create a temporary tunnel (runs in foreground) |
| `tunnerse list` | List all registered tunnels |
| `tunnerse info <tunnel_id>` | Show detailed information about a tunnel |
| `tunnerse kill <tunnel_id>` | Stop a running tunnel |
| `tunnerse del <tunnel_id>` | Delete an inactive tunnel |
| `tunnerse logs <tunnel_id>` | Stream tunnel logs |

## Configuration

The local daemon reads optional `.env` values and has sensible defaults:

| Variable | Default | Description |
| --- | --- | --- |
| `HTTPPort` | `9988` | Port for the local daemon API |
| `SUBDOMAIN` | `false` | Use subdomain routing (true) or path routing (false) |
| `WARNS_ON_HTML` | `true` | Emit HTML warning pages in some failures |
| `TUNNEL_LIFE_TIME` | `86400` | Max lifetime for a tunnel in seconds |
| `TUNNEL_INACTIVITY_LIFE_TIME` | `86400` | Inactivity timeout in seconds |

> By default, the CLI targets `https://tunnerse.com` as the remote API. The local daemon accepts `server_url` in its `/new` and `/quick` endpoints if you want to point to a different API.

## Data & logs

Tunnerse stores local data in:

```
~/.tunnerse/
	├─ db.sqlite
	└─ logs/
```

Each tunnel writes its own log file inside `~/.tunnerse/logs`.

## Linux install helper

There is an optional systemd helper script:

```bash
sudo ./scripts/install.sh
```

This installs `tunnerse` and `tunnerse-server` to `/usr/local/bin` and registers a systemd service.

## Project status

Tunnerse CLI is built for developer workflows: it focuses on fast local exposure with clear logs and simple commands. It is ideal for development and testing environments.

## License

See `LICENSE.md` for details.
