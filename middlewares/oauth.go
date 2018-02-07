package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dictyBase/apihelpers/apherror"
	"github.com/dictyBase/authserver/user"

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

func (m *OauthMiddleware) ParamsMiddleware(h http.Handler) http.Handler {
	if len(m.ConfigParam) == 0 {
		m.ConfigParam = "config"
	}
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		for _, p := range []string{"client_id", "scopes", "redirect_url", "state", "code"} {
			v := r.FormValue(p)
			if len(v) == 0 {
				apherror.JSONAPIError(w,
					apherror.ErrQueryParam.New(
						fmt.Sprintf("missing param %q", p),
					),
				)
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
		newCtx := context.WithValue(ctx, m.ConfigParam, oauthConf)
		h.ServeHTTP(w, r.WithContext(newCtx))
	}
	return http.HandlerFunc(fn)
}

func GetGoogleMiddleware(p *ProvidersSecret) *OauthMiddleware {
	return &OauthMiddleware{
		ClientSecret: p.Google,
		Endpoint:     google.Endpoint,
	}
}

func (m *OauthMiddleware) GoogleMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oauthConf, ok := ctx.Value("config").(*OauthConfig)
		if !ok {
			apherror.JSONAPIError(w, apherror.ErrReqContext.New("no oauth config in request context"))
			return
		}
		oauthConf.Config.ClientSecret = m.ClientSecret
		oauthConf.Config.Endpoint = m.Endpoint
		token, err := oauthConf.Exchange(oauth2.NoContext, oauthConf.Code)
		if err != nil {
			apherror.JSONAPIError(w, apherror.ErrOauthExchange.New(err.Error()))
			return
		}
		oauthClient := oauthConf.Client(oauth2.NoContext, token)
		resp, err := oauthClient.Get(user.Google)
		if err != nil {
			apherror.JSONAPIError(w, apherror.ErrUserRetrieval.New(err.Error()))
			return
		}
		var google user.GoogleUser
		if err := json.NewDecoder(resp.Body).Decode(&google); err != nil {
			apherror.JSONAPIError(w, apherror.ErrJSONEncoding.New(err.Error()))
			return
		}
		user := &user.NormalizedUser{
			Name:  google.Name,
			Email: google.Email,
			Id:    google.Id,
		}
		newCtx := context.WithValue(ctx, "user", user)
		h.ServeHTTP(w, r.WithContext(newCtx))
	}
	return http.HandlerFunc(fn)
}

func GetFacebookMiddleware(p *ProvidersSecret) *OauthMiddleware {
	return &OauthMiddleware{
		ClientSecret: p.Facebook,
		Endpoint:     facebook.Endpoint,
	}
}

func (m *OauthMiddleware) FacebookMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oauthConf, ok := ctx.Value("config").(*OauthConfig)
		if !ok {
			apherror.JSONAPIError(w, apherror.ErrReqContext.New("no oauth config in request context"))
			return
		}
		oauthConf.Config.ClientSecret = m.ClientSecret
		oauthConf.Config.Endpoint = m.Endpoint
		token, err := oauthConf.Exchange(oauth2.NoContext, oauthConf.Code)
		if err != nil {
			apherror.JSONAPIError(w, apherror.ErrOauthExchange.New(err.Error()))
			return
		}
		oauthClient := oauthConf.Client(oauth2.NoContext, token)
		resp, err := oauthClient.Get(user.Facebook)
		if err != nil {
			apherror.JSONAPIError(w, apherror.ErrUserRetrieval.New(err.Error()))
			return
		}
		var facebook user.GoogleUser
		if err := json.NewDecoder(resp.Body).Decode(&facebook); err != nil {
			apherror.JSONAPIError(w, apherror.ErrJSONEncoding.New(err.Error()))
			return
		}
		user := &user.NormalizedUser{
			Name:  facebook.Name,
			Email: facebook.Email,
			Id:    facebook.Id,
		}
		newCtx := context.WithValue(ctx, "user", user)
		h.ServeHTTP(w, r.WithContext(newCtx))
	}
	return http.HandlerFunc(fn)
}

func GetLinkedinMiddleware(p *ProvidersSecret) *OauthMiddleware {
	return &OauthMiddleware{
		ClientSecret: p.LinkedIn,
		Endpoint:     linkedin.Endpoint,
	}
}

func (m *OauthMiddleware) LinkedInMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oauthConf, ok := ctx.Value("config").(*OauthConfig)
		if !ok {
			apherror.JSONAPIError(w, apherror.ErrReqContext.New("no oauth config in request context"))
			return
		}
		oauthConf.Config.ClientSecret = m.ClientSecret
		oauthConf.Config.Endpoint = m.Endpoint
		token, err := oauthConf.Exchange(oauth2.NoContext, oauthConf.Code)
		if err != nil {
			apherror.JSONAPIError(w, apherror.ErrOauthExchange.New(err.Error()))
			return
		}
		oauthClient := oauthConf.Client(oauth2.NoContext, token)
		resp, err := oauthClient.Get(user.LinkedIn)
		if err != nil {
			apherror.JSONAPIError(w, apherror.ErrUserRetrieval.New(err.Error()))
			return
		}
		var linkedin user.LinkedInUser
		if err := json.NewDecoder(resp.Body).Decode(&linkedin); err != nil {
			apherror.JSONAPIError(w, apherror.ErrJSONEncoding.New(err.Error()))
			return
		}
		user := &user.NormalizedUser{
			Name:  fmt.Sprintf("%s %s", linkedin.FirstName, linkedin.LastName),
			Email: linkedin.EmailAddress,
			Id:    linkedin.Id,
		}
		newCtx := context.WithValue(ctx, "user", user)
		h.ServeHTTP(w, r.WithContext(newCtx))
	}
	return http.HandlerFunc(fn)
}
