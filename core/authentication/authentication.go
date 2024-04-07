package authentication

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type ctxKey int

const CtxValues ctxKey = 1

const RoleUser, RoleAdmin = "user", "admin"

type Claims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

type Authentication struct {
	keys map[string]*rsa.PrivateKey
}

func NewAuthentication(keys map[string]*rsa.PrivateKey) *Authentication {
	return &Authentication{keys: keys}
}

func (a *Authentication) GenerateToken(keyId string, claims Claims) (string, error) {
	privateKey, ok := a.keys[keyId]
	if !ok {
		return "", errors.New("invalid key id provided")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = keyId

	return token.SignedString(privateKey)
}

func (a *Authentication) Authenticate(token string) (Claims, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		keyId, ok := token.Header["kid"]
		if !ok {
			return nil, errors.New("key id is missed in token")
		}

		keyIdValue, ok := keyId.(string)
		if !ok {
			return nil, errors.New("invalid key id value provided")
		}

		privateKey, ok := a.keys[keyIdValue]
		if !ok {
			return nil, errors.New("invalid key id provided")
		}

		return &privateKey.PublicKey, nil
	})
	if err != nil {
		return Claims{}, err
	}

	return claims, nil
}

func (a *Authentication) Authorize(claims Claims, roles ...string) bool {
	for _, givenRole := range claims.Roles {
		for _, neededRole := range roles {
			if givenRole == neededRole {
				return true
			}
		}
	}

	return false
}
