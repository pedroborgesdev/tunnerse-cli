package dto

const (
	Welcome string = `Tunnerse - Reverse tunnel manager for development and testing.

`
	Commands string = `Usage: tunnerse <command> [arguments]

Commands:
  new <name> <port>      Create a persistent tunnel (runs in background)
  quick <name> <port>    Create a temporary tunnel (runs in foreground)
  list                   List all registered tunnels
  info <tunnel_id>       Show detailed information about a tunnel
  kill <tunnel_id>       Stop a running tunnel
  del <tunnel_id>        Delete an inactive tunnel from database
  logs <tunnel_id>       View tunnel logs in real-time

Options:
  -h, --help            Show this help message

Examples:
  tunnerse new api-fdp 8080
  tunnerse quick test-app 3000
  tunnerse list
  tunnerse logs api-fdp

`
	Help string = `Tunnerse creates a tunnel that connects the target server using Tunnerse
Server, with your machine pointing to a local port. Tunnerse Server
acts only as an intermediary between the requester and your machine,
while Tunnerse CLI translates the request coming from the server to your
local application. The same process occurs when returning the response
from your application.

Usage: tunnerse <command> [arguments]

Commands:
  new <name> <port>      Create a persistent tunnel (runs in background)
  quick <name> <port>    Create a temporary tunnel (runs in foreground)
  list                   List all registered tunnels
  info <tunnel_id>       Show detailed information about a tunnel
  kill <tunnel_id>       Stop a running tunnel
  del <tunnel_id>        Delete an inactive tunnel from database
  logs <tunnel_id>       View tunnel logs in real-time

Options:
  -h, --help            Show this help message

Examples:
  tunnerse new api-fdp 8080       # Create persistent tunnel
  tunnerse quick test-app 3000   # Create temporary tunnel
  tunnerse list                  # List all tunnels
  tunnerse info api-fdp           # Show tunnel details
  tunnerse kill api-fdp           # Stop tunnel
  tunnerse del api-fdp            # Delete inactive tunnel
  tunnerse logs api-fdp           # View logs

Thanks for using Tunnerse ;)

`

	Start string = `From now on, a tunnel will be created connecting
your local application to the entire internet. thank you
for choosing Tunnerse for this!

Disclaimer:
 Tunnerse only provides the tunnel connection. we are not responsible
 for transmitted content or user data.

`
	Info string = `Version: 1.0.1
Author: pedroborgesdev (on GitHub)
Project: https://github.com/pedroborgesdev/tunnerse-cli.git

`

	Invalid string = `Invalid command usage, use 'help' to see valid usages
	
`

	BetaWarn string = "\033[33mBeta warn:\033[0m \n tunnel ID is passed in the URL path, which may\n" +
		" cause issues with internal navigation links. we attempt to rewrite paths\n" +
		" by placing the ID before the path. this is experimental and may fail.\n\n"

	InvalidID string = `Invalid tunnel id. correct usage should only contain lowercase letters and the special character '-'.
	
`

	InvalidPort string = `Invalid port. correct usage would be only numbers from 0 to 65535.
	
`
)
