package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cyclopsci/apollo"
	"github.com/dictybase/authserver/user"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/linkedin"
)

type ProvidersSecret struct {
	Github   string `json:"github"`
	Facebook string `json:"facebook"`
	Google   string `json:"google"`
	LinkedIn string `json:"linkedin"`
}

type OauthConfig struct {
	State string
	Code  string
	*oauth2.Config
}

type OauthMiddleware struct {
	ClientSecret string
	Endpoint     oauth2.Endpoint
	ConfigParam  string
}

func (m *OauthMiddleware) ParamsMiddleware(h apollo.Handler) apollo.Handler {
	if len(m.ConfigParam) == 0 {
		m.ConfigParam = "config"
	}
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		for _, p := range []string{"client_id", "scopes", "redirect_url", "state", "code"} {
			v := r.FormValue(p)
			if len(v) == 0 {
				http.Error(w, fmt.Sprintf("missing param %q", p), http.StatusBadRequest)
				return
			}
		}
		oauthConf := &OauthConfig{
			Config: &oauth2.Config{
				ClientID:    r.FormValue("client_id"),
				RedirectURL: r.FormValue("redirect_url"),
				Scopes:      strings.Split(r.FormValue("scopes"), " "),
			},
			State: r.FormValue("state"),
			Code:  r.FormValue("code"),
		}
		h.ServeHTTP(context.WithValue(ctx, m.ConfigParam, oauthConf), w, r)
	}
	return apollo.HandlerFunc(fn)
}

func GetGoogleMiddleware(p *ProvidersSecret) *OauthMiddleware {
	return &OauthMiddleware{
		ClientSecret: p.Google,
		Endpoint:     google.Endpoint,
	}
}

func (m *OauthMiddleware) GoogleMiddleware(h apollo.Handler) apollo.Handler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		oauthConf, ok := ctx.Value("config").(*OauthConfig)
		if !ok {
			http.Error(w, "unable to retrieve context", http.StatusInternalServerError)
			return
		}
		oauthConf.Config.ClientSecret = m.ClientSecret
		oauthConf.Config.Endpoint = m.Endpoint
		token, err := oauthConf.Exchange(oauth2.NoContext, oauthConf.Code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		oauthClient := oauthConf.Client(oauth2.NoContext, token)
		resp, err := oauthClient.Get(user.Google)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var google user.GoogleUser
		if err := json.NewDecoder(resp.Body).Decode(&google); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user := &user.NormalizedUser{
			Name:  google.Name,
			Email: google.Email,
		}
		h.ServeHTTP(context.WithValue(ctx, "user", user), w, r)
	}
	return apollo.HandlerFunc(fn)
}

func GetFacebookMiddleware(p *ProvidersSecret) *OauthMiddleware {
	return &OauthMiddleware{
		ClientSecret: p.Facebook,
		Endpoint:     facebook.Endpoint,
	}
}

func (m *OauthMiddleware) FacebookMiddleware(h apollo.Handler) apollo.Handler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		oauthConf, ok := ctx.Value("config").(*OauthConfig)
		if !ok {
			http.Error(w, "unable to retrieve context", http.StatusInternalServerError)
			return
		}
		oauthConf.Config.ClientSecret = m.ClientSecret
		oauthConf.Config.Endpoint = m.Endpoint
		token, err := oauthConf.Exchange(oauth2.NoContext, oauthConf.Code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		oauthClient := oauthConf.Client(oauth2.NoContext, token)
		resp, err := oauthClient.Get(user.Facebook)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var facebook user.GoogleUser
		if err := json.NewDecoder(resp.Body).Decode(&facebook); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user := &user.NormalizedUser{
			Name:  facebook.Name,
			Email: facebook.Email,
		}
		h.ServeHTTP(context.WithValue(ctx, "user", user), w, r)
	}
	return apollo.HandlerFunc(fn)
}

func GetLinkedinMiddleware(p *ProvidersSecret) *OauthMiddleware {
	return &OauthMiddleware{
		ClientSecret: p.LinkedIn,
		Endpoint:     linkedin.Endpoint,
	}
}

func (m *OauthMiddleware) LinkedInMiddleware(h apollo.Handler) apollo.Handler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		oauthConf, ok := ctx.Value("config").(*OauthConfig)
		if !ok {
			http.Error(w, "unable to retrieve context", http.StatusInternalServerError)
			return
		}
		oauthConf.Config.ClientSecret = m.ClientSecret
		oauthConf.Config.Endpoint = m.Endpoint
		token, err := oauthConf.Exchange(oauth2.NoContext, oauthConf.Code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		oauthClient := oauthConf.Client(oauth2.NoContext, token)
		resp, err := oauthClient.Get(user.LinkedIn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var linkedin user.LinkedInUser
		if err := json.NewDecoder(resp.Body).Decode(&linkedin); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user := &user.NormalizedUser{
			Name:  fmt.Sprintf("%q %q", linkedin.FirstName, linkedin.LastName),
			Email: linkedin.EmailAddress,
		}
		h.ServeHTTP(context.WithValue(ctx, "user", user), w, r)
	}
	return apollo.HandlerFunc(fn)
}
