package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/cyclopsci/apollo"
	"github.com/dgrijalva/jwt-go"
	"github.com/dictybase/authserver/handler"
	"github.com/dictybase/authserver/middleware"
	"golang.org/x/net/context"
	"gopkg.in/codegangsta/cli.v1"
)

// The list of providers supported by the server
var DefaultProviders = []string{"google", "facebook"}

func main() {
	app := cli.NewApp()
	app.Name = "authserver"
	app.Usage = "oauth server that provides endpoints for managing authentication"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "runs the auth server",
			Action: runServer,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "config, c",
					Usage:  "Config file(required)",
					EnvVar: "OAUTH_CONFIG",
				},
				cli.StringFlag{
					Name:   "pkey, public-key",
					Usage:  "public key file for verifying jwt",
					EnvVar: "JWT_PUBLIC_KEY",
				},
				cli.StringFlag{
					Name:   "private-key, prkey",
					Usage:  "private key file for signning jwt",
					EnvVar: "JWT_PRIVATE_KEY",
				},
				cli.IntFlag{
					Name:  "port, p",
					Usage: "server port",
					Value: 9999,
				},
			},
		},
		{
			Name:   "generate-keys",
			Usage:  "generate rsa key pairs(public and private keys) in pem format",
			Action: generateKeys,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "private, pr",
					Usage: "output file name for private key",
				},
				cli.StringFlag{
					Name:  "public, pub",
					Usage: "output file name for public key",
				},
			},
		},
	}
	app.Run(os.Args)
}

// Runs the http server
func runServer(c *cli.Context) {
	if !c.IsSet("config") {
		log.Fatalln("config file is not provided")
	}
	if !c.IsSet("public-key") {
		log.Fatalln("public key file is not provided")
	}
	if !c.IsSet("private-key") {
		log.Fatalln("private key file is not provided")
	}

	config, err := readSecretConfig(c)
	if err != nil {
		log.Fatalf("Unable to secret config file %q\n", err)
	}
	jt, err := parseJwtKeys(c)
	if err != nil {
		log.Fatalf("Unable to parse keys %q\n", err)
	}

	mux := http.NewServeMux()
	for _, name := range DefaultProviders {
		switch name {
		case "google":
			mw := middleware.GetGoogleMiddleware(config)
			gchain := apollo.New(mw.ParamsMiddleware, mw.GoogleMiddleware).
				With(context.Background()).
				ThenFunc(jt.JwtHandler)
			mux.Handle("/token/google", gchain)
		case "facebook":
			mw := middleware.GetFacebookMiddleware(config)
			gchain := apollo.New(mw.ParamsMiddleware, mw.FacebookMiddleware).
				With(context.Background()).
				ThenFunc(jt.JwtHandler)
			mux.Handle("/token/facebook", gchain)
		default:
			log.Fatalf("provider %q is not supported\n", name)
		}
	}
	log.Printf("Starting web server on port %d\n", c.Int("port"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Int("port")), mux))

}

// Reads the configuration file containing the various client secret keys
// of the providers. The expected format will be ...
//  {
//		"google": "xxxxxxxxxxxx",
//		"github": "xxxxxxxx",
//	}
func readSecretConfig(c *cli.Context) (*middleware.ProvidersSecret, error) {
	var provider *middleware.ProvidersSecret
	reader, err := os.Open(c.String("config"))
	if err != nil {
		return provider, err
	}
	if err := json.NewDecoder(reader).Decode(&provider); err != nil {
		return provider, err
	}
	return provider, nil

}

// Reads the public and private keys from their respective files and
// stores them in the jwt handler.
func parseJwtKeys(c *cli.Context) (*handler.Jwt, error) {
	jh := &handler.Jwt{}
	private, err := ioutil.ReadFile(c.String("private-key"))
	if err != nil {
		return jh, err
	}
	pkey, err := jwt.ParseRSAPrivateKeyFromPEM(private)
	if err != nil {
		return jh, err
	}

	public, err := ioutil.ReadFile(c.String("public-key"))
	if err != nil {
		return jh, err
	}
	pubkey, err := jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		return jh, err
	}
	jh.VerifyKey = pubkey
	jh.SignKey = pkey
	return jh, err
}

// Generate RSA public and private keys in PEM format
func generateKeys(c *cli.Context) {
	// validate
	if !c.IsSet("public") {
		log.Fatal("public key output file is not provided")
	}
	if !c.IsSet("private") {
		log.Fatal("private key output file is not provided")
	}

	// open files
	prvWriter, err := os.Create(c.String("private"))
	if err != nil {
		log.Fatalf("unable to create private key file %q\n", err)
	}
	pubWriter, err := os.Create(c.String("public"))
	if err != nil {
		log.Fatalf("unable to create public key file %q\n", err)
	}

	// generate and write to files
	private, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("error in generating private key %q\n", err)
	}
	if err := private.Validate(); err != nil {
		log.Fatalf("error in validating private key %q\n", err)
	}
	prvPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(private),
	}
	public := private.Public()
	pubCont, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		log.Fatalf("unable to marshall private key %q\n", err)
	}
	pubPem := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubCont,
	}
	if err := pem.Encode(prvWriter, prvPem); err != nil {
		log.Fatalf("unable to write private key %q\n", err)
	}
	if err := pem.Encode(pubWriter, pubPem); err != nil {
		log.Fatalf("unable to write public key %q\n", err)
	}

}
