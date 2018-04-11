# authserver
An authorization server to generate [jwt](http://jwt.io) token in exchange
for temporary [authorization token](https://tools.ietf.org/html/rfc6749#section-1.4)  
from various [oauth2](http://oauth.net/2/) providres(google, facebook, github, linkedin etc ..).
The server also validate the *jwt* token.

This server is exclusively designed to work with a single page(SPA) frontend web application, for example
something that developed with [React](http://facebook.github.io/react/index.html).

# Supported providers
* [Google](https://developers.google.com/identity/protocols/OAuth2UserAgent)
* [Facebook](https://developers.facebook.com/docs/facebook-login/manually-build-a-login-flow)
* [LinkedIn](https://developer.linkedin.com/docs/oauth2)
* [ORCiD](https://members.orcid.org/api/about-orcid-apis)

# Install
Use the provided `helm`
[chart](https://github.com/dictybase-docker/kubernetes-charts/tree/master/authserver)
to run the server. However, the key generation subcommand can be used
independently. For that either download it from the release
[page](https://github.com/dictyBase/authserver/releases) or install using `go
get`

```
go get github.com/dictyBase/authserver
```

# API
## HTTP/JSON
It's documented [here](https://dictybase.github.io/dictybase-api/), select the `auth` spec from the dropdown.

# Usage
## Generate keys
### Using the subcommand
```
authserver generate-keys --private app.rsa --public app.rsa.pub
```
### Openssl command line(Recommended)
```
openssl genrsa -out keys/app.rsa 2048
openssl rsa -in keys/app.rsa -pubout -out keys/app.rsa.pub 
```
## Create configuration file
The json formatted configuration file should contain `client secret key` for various providers. The secret key
could be obtained by registering a web application with the respective providers.

__Format__

{
    "google": "secret-key-xxxxxxxxxxx",
    "facebook": "secret-key-xxxxxxxxxxx"
    ...........
}

## Command line
```
NAME:
   authserver - oauth server that provides endpoints for managing authentication

USAGE:
   authserver [global options] command [command options] [arguments...]

VERSION:
   4.0.0

COMMANDS:
     run            runs the auth server
     generate-keys  generate rsa key pairs(public and private keys) in pem format
     help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log value, -l value  Name of the log file(optional), default goes to stderr
   --log-format value     Format of the log output,could be either of text or json, default is json
   --help, -h             show help
   --version, -v          print the version
```

### Subcommands
```
NAME:
   authserver run - runs the auth server

USAGE:
   authserver run [command options] [arguments...]

OPTIONS:
   --config value, -c value            Config file(required) [$OAUTH_CONFIG]
   --pkey value, --public-key value    public key file for verifying jwt [$JWT_PUBLIC_KEY]
   --private-key value, --prkey value  private key file for signning jwt [$JWT_PRIVATE_KEY]
   --port value, -p value              server port (default: 9999)
   --messaging-host value              host address for messaging server [$NATS_SERVICE_HOST]
   --messaging-port value              port for messaging server [$NATS_SERVICE_PORT]
```

```
NAME:
   authserver generate-keys - generate rsa key pairs(public and private keys) in pem format

USAGE:
   authserver generate-keys [command options] [arguments...]

OPTIONS:
   --private value, --pr value  output file name for private key
   --public value, --pub value  output file name for public key
```
