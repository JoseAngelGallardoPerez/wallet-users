package auth

import (
	"github.com/Confialink/wallet-pkg-utils"
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/jwt"

	"errors"
	base "github.com/dgrijalva/jwt-go"
	"time"
)

const (
	claimTmpAuthSub = "limited_auth"
	claimTmpAuthExp = "24h"
)

type TemporaryTokens struct {
	jwt jwt.Service
}

func NewTemporaryTokens(jwt jwt.Service) *TemporaryTokens {
	return &TemporaryTokens{jwt: jwt}
}

func (t *TemporaryTokens) Issue(user *models.User) (string, error) {
	duration := utils.MustParseDuration(claimTmpAuthExp)
	exp := time.Now().Add(duration)

	claims := base.MapClaims{
		"sub": claimTmpAuthSub,
		"exp": exp.Unix(),
		"uid": user.UID,
	}

	token := t.jwt.Issue(claims)
	jwtSigned, err := t.jwt.Sign(token)
	if err != nil {
		return "", err
	}
	return jwtSigned, nil
}

func (t *TemporaryTokens) Verify(signedToken string) (*base.Token, error) {
	token, err := t.jwt.Parse(signedToken)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims := token.Claims.(base.MapClaims)

	if claims["sub"] != claimTmpAuthSub {
		return nil, errors.New("invalid token subject")
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, errors.New("token is expired")
	}

	return token, nil
}
