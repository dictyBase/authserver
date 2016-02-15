#authserver
An authorization server to generate [jwt](http://jwt.io) token in exchange
for temporary [authorization token](https://tools.ietf.org/html/rfc6749#section-1.4)  
from various [oauth2](http://oauth.net/2/) providres(google, facebook, github, linkedin etc ..)

This server is exclusively designed to work with a single page(SPA) frontend web application, for example
something that developed with [React](http://facebook.github.io/react/index.html).

#Supported providers
* [Google](https://developers.google.com/identity/protocols/OAuth2UserAgent)
* [Facebook](https://developers.facebook.com/docs/facebook-login/manually-build-a-login-flow)

#Install
```
go get github.com/dictybase/authserver
```

#Usage
##Generate keys

####Using the subcommand

```authserver generate-keys --private app.rsa --public app.rsa.pub```

####Openssl command line

```
openssl genrsa -out keys/app.rsa 2048
openssl rsa -in keys/app.rsa -pubout -out keys/app.rsa.pub 
```

##Create configuration file
The json formatted configuration file should contain `client secret key` for various providers. The secret key
could be obtained by registering a web application with the respective providers.

__Format__

{
    "google": "secret-key-xxxxxxxxxxx",
    "facebook": "secret-key-xxxxxxxxxxx"
    ...........
}


##Run server
```
authserver serve --config app.json --public-key keys/app.rsa.pub --private-key keys/app.rsa
```
The server by default will run in port `9999`

##HTTP post to the server
### Available endpoints
* `/tokens/google` : For google
* `/tokens/facebook` : For facebook

### Required paramater*s
* `client_id` : Available with registered application for every provider.
* `scopes` : Should be available from providers, mostly the value is `email`
* `redirect_url` : As given in the registered application
* `state` : As passed to the provider during the first login
* `code` : As passed to the redirect_url from the provider

An example of http post using `curl`

First write all paramaters to a file, say `params.txt`. The content of the file will look like
```
client_id=xxxxxxx&scopes=email&redirect_url=http://localhost:3000/google/callback&state=google&code=xxxxxx
```

```
curl -X POST -d @params.txt http://localhost:9999/tokens/google
```
The above should a return a `json web token`.

##Command line
```
NAME:
   authserver - oauth server that provides endpoints for managing authentication

USAGE:
   authserver [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   run			runs the auth server
   generate-keys	generate rsa key pairs(public and private keys) in pem format
   help, h		Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --log, -l            Name of the log file(optional), default goes to stderr
   --help, -h		show help
   --version, -v	print the version
   
```

#### Subcommands
```
NAME:
   generate-keys - generate rsa key pairs(public and private keys) in pem format

USAGE:
   command generate-keys [command options] [arguments...]

DESCRIPTION:
   

OPTIONS:
   --private, --pr 	output file name for private key
   --public, --pub 	output file name for public key
``` 
```
NAME:
   run - runs the auth server

USAGE:
   command run [command options] [arguments...]

DESCRIPTION:
   

OPTIONS:
   --config, -c 		Config file(required) [$OAUTH_CONFIG]
   --pkey, --public-key 	public key file for verifying jwt [$JWT_PUBLIC_KEY]
   --private-key, --prkey 	private key file for signning jwt [$JWT_PRIVATE_KEY]
   --port, -p '9999'		server port
```

