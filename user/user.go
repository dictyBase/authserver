// package user provides url constants and structures
// for handling user information
// from various oauth providers
package user // import "github.com/dictybase/authserver/user"

var (
	Google   = "https://www.googleapis.com/userinfo/v2/me"
	Facebook = "https://graph.facebook.com/v2.5/me?fields=name,email"
)

type GoogleUser struct {
	FamilyName    string `json:"family_name"`
	Name          string `json:"name"`
	Gender        string `json:"gender"`
	Email         string `json:"email"`
	GivenName     string `json:"given_name"`
	ID            string `json:"id"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

type FacebookUser struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type NormalizedUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
