package jwt

import (
	base "github.com/dgrijalva/jwt-go"
)

type WithHMAC struct {
	signingMethod *base.SigningMethodHMAC
	secret        string
}

func NewWithHMAC(secret string, signingMethod *base.SigningMethodHMAC) (service Service) {
	service = &WithHMAC{secret: secret, signingMethod: signingMethod}
	return
}

func (s *WithHMAC) Issue(claims base.Claims) *base.Token {
	return base.NewWithClaims(s.signingMethod, claims)
}

func (s *WithHMAC) Sign(t *base.Token, secret ...[]byte) (string, error) {
	key := []byte(s.secret)
	if len(secret) > 0 {
		key = secret[0]
	}
	return t.SignedString(key)
}

func (s *WithHMAC) Parse(str string, secret ...[]byte) (*base.Token, error) {
	key := []byte(s.secret)
	if len(secret) > 0 {
		key = secret[0]
	}
	return base.Parse(str, func(token *base.Token) (interface{}, error) {
		return key, nil
	})
}
