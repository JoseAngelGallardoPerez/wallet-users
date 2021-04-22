package jwt

import (
	base "github.com/dgrijalva/jwt-go"
)

const DefaultSecretKeyPath = "./jwt.pem"
const DefaultPublicKeyPath = "./jwt.pub"

type Issuer interface {
	Issue(claims base.Claims) *base.Token
}

type Signer interface {
	Sign(t *base.Token, secret ...[]byte) (string, error)
}

type Parser interface {
	Parse(s string, secret ...[]byte) (*base.Token, error)
}

type Service interface {
	Issuer
	Signer
	Parser
}
