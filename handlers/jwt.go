package handlers

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dictybase/authserver/user"
	"golang.org/x/net/context"
)

type Jwt struct {
	VerifyKey     *rsa.PublicKey
	SignKey       *rsa.PrivateKey
	UserParamater string
}

type UserToken struct {
	Token string `json:"token"`
}

func (j *Jwt) JwtHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// default value for user paramater
	if len(j.UserParamater) == 0 {
		j.UserParamater = "user"
	}
	user, ok := ctx.Value(j.UserParamater).(*user.NormalizedUser)
	if !ok {
		http.Error(w, "unable to retrieve user from context", http.StatusInternalServerError)
		return
	}
	// create a signer for rsa 512
	t := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"iss":  "dictyBase",
		"sub":  "dictyBase login token",
		"exp":  time.Now().Add(time.Hour * 240).Unix(),
		"iat":  time.Now().Unix(),
		"user": user, // custom claim
	})
	token, err := t.SignedString(j.SignKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(&UserToken{token}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
