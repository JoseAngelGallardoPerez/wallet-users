package jwt

import (
	"crypto/ecdsa"
	base "github.com/dgrijalva/jwt-go"
	"io/ioutil"
)

type WithECDSA struct {
	signingMethod *base.SigningMethodECDSA
	secretKey     *ecdsa.PrivateKey
	publicKey     *ecdsa.PublicKey
}

type Params struct {
	PublicKey     []byte
	SecretKey     []byte
	SigningMethod *base.SigningMethodECDSA
	SecretKeyPath string
	PublicKeyPath string
}

func NewWithECDSA(params *Params) (service Service, err error) {
	if params == nil {
		params = &Params{
			SecretKeyPath: DefaultSecretKeyPath,
			PublicKeyPath: DefaultPublicKeyPath,
		}
	}

	var secretBytes []byte
	if params.SecretKey == nil {
		secretBytes, err = ioutil.ReadFile(params.SecretKeyPath)
		if nil != err {
			return
		}
		params.SecretKey = secretBytes
	}
	var publicBytes []byte
	if params.PublicKey == nil {
		publicBytes, err = ioutil.ReadFile(params.PublicKeyPath)
		if nil != err {
			return
		}
		params.PublicKey = publicBytes
	}

	secretKey, err := base.ParseECPrivateKeyFromPEM(params.SecretKey)
	if nil != err {
		return
	}

	publicKey, err := base.ParseECPublicKeyFromPEM(params.PublicKey)
	if nil != err {
		return
	}

	if params.SigningMethod == nil {
		params.SigningMethod = base.SigningMethodES512
	}

	service = &WithECDSA{signingMethod: params.SigningMethod, secretKey: secretKey, publicKey: publicKey}
	return
}

func (s *WithECDSA) Issue(claims base.Claims) *base.Token {
	return base.NewWithClaims(s.signingMethod, claims)
}

func (s *WithECDSA) Sign(t *base.Token, secret ...[]byte) (string, error) {
	return t.SignedString(s.secretKey)
}

func (s *WithECDSA) Parse(str string, secret ...[]byte) (*base.Token, error) {
	return base.Parse(str, func(token *base.Token) (interface{}, error) {
		return s.publicKey, nil
	})
}
