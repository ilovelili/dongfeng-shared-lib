package sharedlib

import (
	"errors"
	"fmt"
	"sync"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kelvinji2009/graphql"
	"github.com/lestrrat-go/jwx/jwk"

	authing "github.com/Authing/authing-go-sdk"
	prettyjson "github.com/hokaccha/go-prettyjson"
)

var (
	keyset   *jwk.Set
	once     sync.Once
	instance *Client
)

// Client authing client wrapper
type Client struct {
	AuthingClient *authing.Client
}

// NewAuthClient  new authing client
func NewAuthClient(clientID, appSecret string) *Client {
	once.Do(func() {
		client := authing.NewClient(clientID, appSecret, false)
		// Enable debug info for graphql client, just comment it if you want to disable the debug info
		client.Client.Log = func(s string) {
			b := []byte(s)
			pj, _ := prettyjson.Format(b)
			fmt.Println(string(pj))
		}

		instance = &Client{
			AuthingClient: client,
		}
	})

	return instance
}

// VerifyLogin verifiy login status by idtoken
func (c *Client) VerifyLogin(idtoken string) (bool, error) {
	p := authing.CheckLoginStatusQueryParameter{
		Token: graphql.String(idtoken),
	}

	q, err := c.AuthingClient.CheckLoginStatus(&p)
	if err != nil {
		return false, err
	}

	return bool(q.CheckLoginStatus.Status), nil
}

// ParseJWT parse JWT
func ParseJWT(idtoken, jwks string) (claims jwt.MapClaims, token *jwt.Token, err error) {
	// Get json web token key set https://auth0.com/docs/jwks
	ks, err := fetchKeySet(jwks)
	if err != nil {
		return
	}

	keyset = ks
	claims = jwt.MapClaims{}
	token, err = jwt.ParseWithClaims(idtoken, claims, getKey)
	return
}

func getKey(token *jwt.Token) (interface{}, error) {
	// validate the alg
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("expecting JWT header to have string kid")
	}

	if key := keyset.LookupKeyID(keyID); len(key) == 1 {
		return key[0].Materialize()
	}

	return nil, errors.New("unable to find key")
}

// fetchKeySet fetch key set
func fetchKeySet(jwks string) (*jwk.Set, error) {
	return jwk.FetchHTTP(jwks)
}
