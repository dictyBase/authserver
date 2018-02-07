// package user provides url constants and structures
// for handling user information
// from various oauth providers
package user // import "github.com/dictybase/authserver/user"

var (
	Google   = "https://www.googleapis.com/userinfo/v2/me"
	Facebook = "https://graph.facebook.com/v2.5/me?fields=name,email"
	LinkedIn = "https://api.linkedin.com/v1/people/~:(first-name,last-name,email-address)?format=json"
)

type GoogleUser struct {
	FamilyName    string `json:"family_name"`
	Name          string `json:"name"`
	Gender        string `json:"gender"`
	Email         string `json:"email"`
	GivenName     string `json:"given_name"`
	Id            string `json:"id"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

type FacebookUser struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LinkedInUser struct {
	FirstName                  string `json:"firstName"`
	Headline                   string `json:"headline"`
	Id                         string `json:"id"`
	LastName                   string `json:"lastName"`
	SiteStandardProfileRequest struct {
		URL string `json:"url"`
	} `json:"siteStandardProfileRequest"`
	EmailAddress string `json:"emailAddress"`
}

type NormalizedUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Id    string `json:"id"`
}
