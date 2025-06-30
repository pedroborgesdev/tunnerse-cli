package dto

const (
	Welcome string = `Tunnerse - temporary reverse tunnel for development and testing.

`
	Commands string = `usage: tunnerse <command>

commands:
 <tunnel_name> <local_port>    create a tunnel that exposes your
                               local application to the web.

 verison    see the version of the CLI.
 help       show this.
	
`
	Help string = `Tunnerse creates a tunnel that connects the target server using Tunnerse
Server, with your machine pointing to a local port. Tunnerse Server
acts only as an intermediary between the requester and your machine,
while Tunnerse CLI translates the request coming from the server to your
local application. the same process occurs when returning the response
from your application.

usage: tunnerse <command>

commands:
 <tunnel_name> <local_port>    create a tunnel that exposes your
                               local application to the web.

 verison    see the version of the CLI.
 help       show this.

thanks for using Tunnerse ;)

`

	Start string = `from now on, a tunnel will be created connecting
your local application to the entire internet. thank you
for choosing Tunnerse for this!

disclaimer:
 tunnerse only provides the tunnel connection. we are not responsible
 for transmitted content or user data.

`
	Info string = `version: 1.0.1
author: pedroborgesdev (on GitHub)
project: https://github.com/pedroborgesdev/tunnerse-cli.git

`

	Invalid string = `invalid command usage, use 'help' to see valid usages
	
`

	BetaWarn string = "\033[33mbeta warn:\033[0m \n tunnel ID is passed in the URL path, which may\n" +
		" cause issues with internal navigation links. we attempt to rewrite paths\n" +
		" by placing the ID before the path. this is experimental and may fail.\n\n"

	InvalidID string = `invalid tunnel id. correct usage should only contain lowercase letters and the special character '-'.
	
`

	InvalidPort string = `invalid port. correct usage would be only numbers from 0 to 65535.
	
`
)
