package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/urfave/cli.v1"

	"github.com/dgrijalva/jwt-go"
	"github.com/dictyBase/authserver/handlers"
	"github.com/dictyBase/authserver/message/nats"
	"github.com/dictyBase/authserver/middlewares"
	"github.com/dictyBase/go-middlewares/middlewares/logrus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
)

// Runs the http server
func RunServer(c *cli.Context) error {
	reqm, err := nats.NewRequest(
		c.String("nats-host"),
		c.String("nats-port"),
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("cannot connect to messaging server %s", err.Error()),
			2,
		)
	}
	config, err := readSecretConfig(c)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Unable to secret config file %q\n", err), 2)
	}
	jt, err := parseJwtKeys(c)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Unable to parse keys %q\n", err), 2)
	}
	// sets the reply messaging connection
	jt.Request = reqm
	jt.Topics = map[string]string{
		"userExists":     "UserService.Exist",
		"userGet":        "UserService.Get",
		"identityExists": "IdentityService.Exist",
		"identityGet":    "IdentityService.Get",
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
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler)
	// Default health check
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("okay"))
	})
	googleMw := middlewares.GetGoogleMiddleware(config)
	fbookMw := middlewares.GetFacebookMiddleware(config)
	linkedInMw := middlewares.GetLinkedinMiddleware(config)
	OrcidMw := middlewares.GetOrcidMiddleware(config)
	r.Route("/tokens/google", func(r chi.Router) {
		r.Use(googleMw.ParamsMiddleware, googleMw.GoogleMiddleware)
		r.Post("/", jt.JwtHandler)
	})

	r.Route("/tokens/facebook", func(r chi.Router) {
		r.With(fbookMw.ParamsMiddleware).
			With(fbookMw.FacebookMiddleware).
			Post("/", jt.JwtHandler)
	})

	r.Route("/tokens/linkedin", func(r chi.Router) {
		r.With(linkedInMw.ParamsMiddleware).
			With(linkedInMw.LinkedInMiddleware).
			Post("/", jt.JwtHandler)
	})

	r.Route("/tokens/orcid", func(r chi.Router) {
		r.With(OrcidMw.ParamsMiddleware).
			With(OrcidMw.OrcidMiddleware).
			Post("/", jt.JwtHandler)
	})

	r.Route("/tokens/validate", func(r chi.Router) {
		tokenAuth := jwtauth.New("RS512", jt.SignKey, jt.VerifyKey)
		r.With(jwtauth.Verifier(tokenAuth), jwtauth.Authenticator).
			Get("/", jt.JwtFinalHandler)
	})
	log.Printf("Starting web server on port %d\n", c.Int("port"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Int("port")), r))
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

// Gets a logrus logger middlware, can output to a file, default is stderr
func getLoggerMiddleware(c *cli.Context) (*logrus.Logger, error) {
	var logger *logrus.Logger
	if c.GlobalIsSet("log") {
		w, err := os.Open(c.GlobalString("log"))
		if err != nil {
			return logger, fmt.Errorf("could not open log file for writing %s\n", err)
		}
		if c.GlobalString("log-format") == "text" {
			logger = logrus.NewFileLogger(w)
		} else {
			logger = logrus.NewJSONFileLogger(w)
		}
	} else {
		if c.GlobalString("log-format") == "text" {
			logger = logrus.NewLogger()
		} else {
			logger = logrus.NewJSONLogger()
		}
	}
	return logger, nil
}
