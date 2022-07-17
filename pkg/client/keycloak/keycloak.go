package keycloak

import (
	"chatty/config"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

const (
	authorizationHeader    = "Authorization"
	formContentType        = "application/x-www-form-urlencoded"
	clientCredentialsGrant = "client_credentials"

	tokenPath     = "/auth/realms/chatty-realm/protocol/openid-connect/token"
	listUsersPath = `/auth/admin/realms/chatty-realm/users`
)

type Keycloak struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Keycloak {
	return &Keycloak{
		cfg: cfg,
	}
}

type Token struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
	Scope            string `json:"scope"`
}

type UserAttributes struct {
	Picture []string `json:"picture"`
}
type User struct {
	ID         uuid.UUID      `json:"id"`
	Created    int            `json:"createdTimestamp"`
	FirstName  string         `json:"firstName"`
	LastName   string         `json:"lastName"`
	Email      string         `json:"email"`
	Username   string         `json:"username"`
	Attributes UserAttributes `json:"attributes"`
}

type UserList []User

func (k *Keycloak) GetToken() (*Token, error) {
	reqUrl := strings.Join([]string{k.cfg.Auth.Base, tokenPath}, "")

	data := url.Values{}
	data.Add("client_id", k.cfg.Auth.ClientID)
	data.Add("client_secret", k.cfg.Auth.Secret)
	data.Add("grant_type", clientCredentialsGrant)

	resp, err := http.Post(reqUrl, formContentType, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	var token Token
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (k *Keycloak) ListUsers(firstName, lastName, email, search string) (UserList, error) {
	token, err := k.GetToken()
	if err != nil {
		return nil, err
	}

	reqUrl := strings.Join([]string{k.cfg.Auth.Base, listUsersPath}, "")

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()

	if firstName != "" {
		q.Add("firstName", firstName)
	}
	if lastName != "" {
		q.Add("lastName", lastName)
	}
	if email != "" {
		q.Add("email", email)
	}
	if search != "" {
		q.Add("search", search)
	}
	req.URL.RawQuery = q.Encode()

	auth := strings.Join([]string{"Bearer", token.AccessToken}, " ")
	req.Header.Set(authorizationHeader, auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var userListResponse UserList
	// var asd interface{}
	err = json.NewDecoder(resp.Body).Decode(&userListResponse)
	if err != nil {
		return nil, err
	}

	return userListResponse, nil
}
