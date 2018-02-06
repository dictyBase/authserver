package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/cyclopsci/apollo"
	"github.com/dgrijalva/jwt-go"
	"github.com/dictyBase/go-middlewares/middlewares/logrus"
	"github.com/dictybase/authserver/handlers"
	"github.com/dictybase/authserver/middlewares"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"golang.org/x/net/context"
	"gopkg.in/codegangsta/cli.v1"
)

// The list of providers supported by the server
var DefaultProviders = []string{"google", "facebook", "linkedin"}

// Runs the http server
func RunServer(c *cli.Context) error {
	if !c.IsSet("config") {
		if len(os.Getenv("OAUTH_CONFIG")) == 0 {
			return cli.NewExitError("config file is not provided", 2)
		}
	}
	if !c.IsSet("public-key") {
		if len(os.Getenv("JWT_PUBLIC_KEY")) == 0 {
			return cli.NewExitError("public key file is not provided", 2)
		}
	}
	if !c.IsSet("private-key") {
		if len(os.Getenv("JWT_PRIVATE_KEY")) == 0 {
			return cli.NewExitError("private key file is not provided", 2)
		}
	}

	config, err := readSecretConfig(c)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Unable to secret config file %q\n", err), 2)
	}
	jt, err := parseJwtKeys(c)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Unable to parse keys %q\n", err), 2)
	}

	loggerMw, err := getLoggerMiddleware(c)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("unable to get logger middlware %s", err), 2)
	}
	cors := cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
	})

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(loggerMw.Middleware)
	r.User(middleware.Recoverer)
	for _, name := range DefaultProviders {
		switch name {
		case "google":
			googleMw := middlewares.GetGoogleMiddleware(config)
			gchain := apollo.New(
				apollo.Wrap(cors.Handler),
				apollo.Wrap(logMw.LoggerMiddleware),
				googleMw.ParamsMiddleware,
				googleMw.GoogleMiddleware).
				With(context.Background()).
				ThenFunc(jt.JwtHandler)
			mux.Handle("/tokens/google", gchain)
		case "facebook":
			fbookMw := middlewares.GetFacebookMiddleware(config)
			fchain := apollo.New(
				apollo.Wrap(cors.Handler),
				apollo.Wrap(logMw.LoggerMiddleware),
				fbookMw.ParamsMiddleware,
				fbookMw.FacebookMiddleware).
				With(context.Background()).
				ThenFunc(jt.JwtHandler)
			mux.Handle("/tokens/facebook", fchain)
		case "linkedin":
			linkeinMw := middlewares.GetLinkedinMiddleware(config)
			lchain := apollo.New(
				apollo.Wrap(cors.Handler),
				apollo.Wrap(logMw.LoggerMiddleware),
				linkeinMw.ParamsMiddleware,
				linkeinMw.LinkedInMiddleware).
				With(context.Background()).
				ThenFunc(jt.JwtHandler)
			mux.Handle("/tokens/linkedin", lchain)
		default:
			return cli.NewExitError(fmt.Sprintf("provider %q is not supported\n", name), 2)
		}
	}
	log.Printf("Starting web server on port %d\n", c.Int("port"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Int("port")), mux))
	return nil
}

// Reads the configuration file containing the various client secret keys
// of the providers. The expected format will be ...
//  {
//		"google": "xxxxxxxxxxxx",
//		"github": "xxxxxxxx",
//	}
func readSecretConfig(c *cli.Context) (*middlewares.ProvidersSecret, error) {
	var provider *middlewares.ProvidersSecret
	reader, err := os.Open(c.String("config"))
	defer reader.Close()
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
func parseJwtKeys(c *cli.Context) (*handlers.Jwt, error) {
	jh := &handlers.Jwt{}
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

func getLoggerMiddleware(c *cli.Context) (*logrus.Logger, error) {
	var logger *logrus.Logger
	if c.GlobalIsSet("log") {
		w, err := os.Open(c.GlobalString("log"))
		if err != nil {
			return logger, fmt.Errorf("could not open log file for writing %s\n", err)
		}
		if c.GlobalString("log-format") == "json" {
			logger = logrus.NewJSONFileLogger(w)
		} else {
			logger = logrus.NewFileLogger(w)
		}
	} else {
		if c.GlobalString("log-format") == "json" {
			logger = logrus.NewJSONLogger()
		} else {
			logger = logrus.NewLogger()
		}
	}
	return logger, nil
}
