package orcid

import (
	"golang.org/x/oauth2"
)

var Endpoint = oauth2.Endpoint{
	AuthURL:  "https://orcid.org/oauth/authorize",
	TokenURL: "https://orcid.org/oauth/token",
}
