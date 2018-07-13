package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	loggerMw "github.com/dictyBase/go-middlewares/middlewares/logrus"
	"gopkg.in/urfave/cli.v1"

	"github.com/dgrijalva/jwt-go"
	"github.com/dictyBase/authserver/handlers"
	"github.com/dictyBase/authserver/message/nats"
	"github.com/dictyBase/authserver/middlewares"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	gnats "github.com/nats-io/go-nats"
)

// Runs the http server
func RunServer(c *cli.Context) error {
	reqm, err := nats.NewRequest(
		c.String("messaging-host"),
		c.String("messaging-port"),
		gnats.MaxReconnects(-1),
		gnats.ReconnectWait(2*time.Second),
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
		"identityGet":    "IdentityService.GetIdentity",
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
		if !reqm.IsActive() {
			http.Error(w, "messaging server is disconnected", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("okay"))
	})
	googleMw := middlewares.GetGoogleMiddleware(config)
	fbookMw := middlewares.GetFacebookMiddleware(config)
	linkedInMw := middlewares.GetLinkedinMiddleware(config)
	OrcidMw := middlewares.GetOrcidMiddleware(config)
	r.Route("/tokens", func(r chi.Router) {
		r.With(googleMw.ParamsMiddleware).
			With(googleMw.GoogleMiddleware).Post("/google", jt.JwtHandler)
		r.With(fbookMw.ParamsMiddleware).
			With(fbookMw.FacebookMiddleware).Post("/facebook", jt.JwtHandler)
		r.With(linkedInMw.ParamsMiddleware).
			With(linkedInMw.LinkedInMiddleware).Post("/linkedin", jt.JwtHandler)
		r.With(OrcidMw.ParamsMiddleware).
			With(OrcidMw.OrcidMiddleware).Post("/orcid", jt.JwtHandler)
	})
	r.Route("/authorize", func(r chi.Router) {
		tokenAuth := jwtauth.New("RS512", jt.SignKey, jt.VerifyKey)
		r.With(middlewares.AuthorizeMiddleware).
			With(jwtauth.Verifier(tokenAuth), jwtauth.Authenticator).
			Post("/", jt.JwtFinalHandler)
	})
	if err := chi.Walk(r, walkFunc); err != nil {
		log.Printf("error in printing routes %s\n", err)
	}
	log.Printf("Starting web server on port %d\n", c.Int("port"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Int("port")), r))
	return nil
}

// Prints all the registered routes
func walkFunc(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
	log.Printf("method: %s - - route: %s\n", method, route)
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

// GetLoggerMiddleware gets a net/http compatible instance of logrus
func getLoggerMiddleware(c *cli.Context) (*loggerMw.Logger, error) {
	var logger *loggerMw.Logger
	var w io.Writer
	if c.IsSet("log-file") {
		fw, err := os.Create(c.String("log-file"))
		if err != nil {
			return logger,
				fmt.Errorf("could not open log file  %s %s", c.String("log-file"), err)
		}
		w = io.MultiWriter(fw, os.Stderr)
	} else {
		w = os.Stderr
	}
	if c.String("log-format") == "json" {
		logger = loggerMw.NewJSONFileLogger(w)
	} else {
		logger = loggerMw.NewFileLogger(w)
	}
	return logger, nil
}
