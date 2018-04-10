package handlers

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc/status"

	"github.com/dictyBase/go-genproto/dictybaseapis/pubsub"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dictyBase/apihelpers/apherror"
	"github.com/dictyBase/authserver/message"
	"github.com/dictyBase/authserver/user"
	"github.com/go-chi/jwtauth"
	"github.com/rs/xid"
)

type contextKey string

// String output the details of context key
func (c contextKey) String() string {
	return "pagination context key " + string(c)
}

var (
	ContextKeyUser = contextKey("user")
)

type Jwt struct {
	VerifyKey     *rsa.PublicKey
	SignKey       *rsa.PrivateKey
	UserParamater string
	Request       message.Request
	Topics        map[string]string
}

type UserToken struct {
	Token string               `json:"token"`
	User  *user.NormalizedUser `json:"user"`
}

func (j *Jwt) JwtFinalHandler(w http.ResponseWriter, r *http.Request) {
	_, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		apherror.JSONAPIError(w, apherror.ErrReqContext.New("unable to retrieve jwt from context for validation"))
		return
	}
	fmt.Fprintf(w, "jwt is %s", "valid")
}

func (j *Jwt) JwtHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(user.ContextKeyUser).(*user.NormalizedUser)
	if !ok {
		apherror.JSONAPIError(w, apherror.ErrReqContext.New("unable to retrieve %s from context", "user"))
		return
	}
	idnReq := &pubsub.IdentityReq{Provider: user.Provider, Identifier: user.Email}
	if user.Provider == "orcid" {
		idnReq.Identifier = user.Id
	}
	// check if the identity is present
	idnReply, err := j.Request.IdentityRequestWithContext(
		context.Background(),
		j.Topics["identityExists"],
		idnReq,
	)

	claims := jwt.StandardClaims{
		Issuer:    "dictyBase",
		Subject:   "dictyBase login token",
		ExpiresAt: time.Now().Add(time.Hour * 240).Unix(),
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
		Id:        xid.New().String(),
		Audience:  "user",
	}

	// create a signer for rsa 512
	t := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	token, err := t.SignedString(j.SignKey)
	if err != nil {
		apherror.ErrJWTToken.New("error in signing jwt token %s", err.Error())
		return
	}
	ut := &UserToken{
		Token: token,
		User:  user,
	}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	if err := json.NewEncoder(w).Encode(ut); err != nil {
		apherror.ErrJSONEncoding.New(err.Error())
		return
	}
}

func handleUserErr(w http.ResponseWriter, reply pubsub.UserReply, id string, err error) bool {
	if err != nil {
		apherror.JSONAPIError(w, apherror.ErrMessagingReply.New("error in getting user reply %s", err.Error()))
		return true
	}
	if reply.Status != nil {
		apherror.JSONAPIError(w, apherror.ErrMessagingReply.New(status.ErrorProto(reply.Status).Error()))
		return true
	}
	if !reply.Exist {
		apherror.JSONAPIError(w, apherror.ErrMessagingReply.New("dictybase user %s not found", id))
		return true
	}
	return false
}

func handleIdentityErr(w http.ResponseWriter, reply pubsub.IdentityReply, id string, err error) bool {
	if err != nil {
		apherror.JSONAPIError(w, apherror.ErrMessagingReply.New("error in getting identifier reply %s", err.Error()))
		return true
	}
	if reply.Status != nil {
		apherror.JSONAPIError(w, apherror.ErrMessagingReply.New(status.ErrorProto(reply.Status).Error()))
		return true
	}
	if !reply.Exist {
		apherror.JSONAPIError(w, apherror.ErrMessagingReply.New("identifier %s not found", id))
		return true
	}
	return false
}
