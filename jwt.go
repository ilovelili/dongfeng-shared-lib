package sharedlib

import (
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
)

var keyset *jwk.Set

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
