// package user provides url constants and structures
// for handling user information
// from various oauth providers
package user

var (
	Google   = "https://www.googleapis.com/userinfo/v2/me"
	Facebook = "https://graph.facebook.com/v2.5/me?fields=name,email"
	LinkedIn = "https://api.linkedin.com/v1/people/~:(first-name,last-name,email-address)?format=json"
	Orcid    = "https://pub.orcid.org/v2.1"
)

type contextKey string

// String output the details of context key
func (c contextKey) String() string {
	return "pagination context key " + string(c)
}

var (
	ContextKeyConfig = contextKey("config")
	ContextKeyUser   = contextKey("user")
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

type OrcidUser struct {
	LastModifiedDate interface{} `json:"last-modified-date"`
	Name             struct {
		CreatedDate struct {
			Value int64 `json:"value"`
		} `json:"created-date"`
		LastModifiedDate struct {
			Value int64 `json:"value"`
		} `json:"last-modified-date"`
		GivenNames struct {
			Value string `json:"value"`
		} `json:"given-names"`
		FamilyName struct {
			Value string `json:"value"`
		} `json:"family-name"`
		CreditName interface{} `json:"credit-name"`
		Source     interface{} `json:"source"`
		Visibility string      `json:"visibility"`
		Path       string      `json:"path"`
	} `json:"name"`
	OtherNames struct {
		LastModifiedDate interface{}   `json:"last-modified-date"`
		OtherName        []interface{} `json:"other-name"`
		Path             string        `json:"path"`
	} `json:"other-names"`
	Biography      interface{} `json:"biography"`
	ResearcherUrls struct {
		LastModifiedDate interface{}   `json:"last-modified-date"`
		ResearcherURL    []interface{} `json:"researcher-url"`
		Path             string        `json:"path"`
	} `json:"researcher-urls"`
	Emails struct {
		LastModifiedDate interface{}   `json:"last-modified-date"`
		Email            []interface{} `json:"email"`
		Path             string        `json:"path"`
	} `json:"emails"`
	Addresses struct {
		LastModifiedDate interface{}   `json:"last-modified-date"`
		Address          []interface{} `json:"address"`
		Path             string        `json:"path"`
	} `json:"addresses"`
	Keywords struct {
		LastModifiedDate interface{}   `json:"last-modified-date"`
		Keyword          []interface{} `json:"keyword"`
		Path             string        `json:"path"`
	} `json:"keywords"`
	ExternalIdentifiers struct {
		LastModifiedDate   interface{}   `json:"last-modified-date"`
		ExternalIdentifier []interface{} `json:"external-identifier"`
		Path               string        `json:"path"`
	} `json:"external-identifiers"`
	Path string `json:"path"`
}

type NormalizedUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Id    string `json:"id"`
}
