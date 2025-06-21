package dto

const (
	Start string = `Thanks for using Tunnerse! ;)
↳ temporary reverse tunnel for development and testing.

version: 1.0.1
author: @pedroborgezs
project: https://github.com/pedroborgzes/tunnerse-cli.git

disclaimer:
 tunnerse only provides the tunnel connection. We are not responsible
 for transmitted content or user data.

`

	Usage string = `usage: tunnerse <tunnel_name> <port>`

	BetaWarn string = "\033[33mbeta warn:\033[0m \n tunnel ID is passed in the URL path, which may\n" +
		" cause issues with internal navigation links. We attempt to rewrite paths\n" +
		" by placing the ID before the path. This is experimental and may fail.\n\n"
)
