package auth

import (
	"time"

	"errors"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/jwt"
	base "github.com/dgrijalva/jwt-go"
)

const (
	ClaimAccessSub       = "access"
	ClaimRefreshSub      = "refresh"
	ClaimAccessTokenExp  = "30m"
	ClaimRefreshTokenExp = "720h"
)

type TokensResponse struct {
	Access  string `json:"accessToken"`
	Refresh string `json:"refreshToken"`
}

type ExtendedTokensResponse struct {
	TokensResponse
	ChallengeName *string `json:"challengeName"`
}

type TokenService struct {
	jwt              jwt.Service
	tokenRepository  *repositories.TokenRepository
	tokenTTLResolver TokenTTLResolver
}

// TokenTTLResolver defines time to life for token
type TokenTTLResolver interface {
	// ResolveByTokenSubject determines time to life by given subject e.g. access, refresh
	ResolveByTokenSubject(subject string) (time.Duration, error)
}

// TokenOptions a set of options which impacts on token fields
type TokenOptions struct {
	// TtlResolver overrides default token resolver if set
	TtlResolver TokenTTLResolver
}

func NewTokenService(jwt jwt.Service, tokenRepository *repositories.TokenRepository, tokenTTLResolver TokenTTLResolver) *TokenService {
	return &TokenService{
		jwt:              jwt,
		tokenRepository:  tokenRepository,
		tokenTTLResolver: tokenTTLResolver,
	}
}

func (t *TokenService) IssueTokens(user *models.User, options *TokenOptions) (*TokensResponse, error) {
	refresh, err := t.IssueRefreshToken(user, options)
	if err != nil {
		return nil, err
	}

	access, err := t.IssueAccessToken(user, refresh, options)
	if err != nil {
		return nil, err
	}

	return &TokensResponse{
		Access:  access.SignedString,
		Refresh: refresh.SignedString,
	}, nil
}

func (t *TokenService) IssueAccessToken(user *models.User, refreshToken *models.Token, options *TokenOptions) (*models.Token, error) {
	var resolver TokenTTLResolver
	if options != nil && options.TtlResolver != nil {
		resolver = options.TtlResolver
	} else {
		resolver = t.tokenTTLResolver
	}
	ttl, err := resolver.ResolveByTokenSubject(ClaimAccessSub)
	if err != nil {
		return nil, err
	}
	model, err := t.issueToken(user, ClaimAccessSub, ttl, &refreshToken.ID)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (t *TokenService) IssueRefreshToken(user *models.User, options *TokenOptions) (*models.Token, error) {
	var resolver TokenTTLResolver
	if options != nil && options.TtlResolver != nil {
		resolver = options.TtlResolver
	} else {
		resolver = t.tokenTTLResolver
	}
	ttl, err := resolver.ResolveByTokenSubject(ClaimRefreshSub)
	if err != nil {
		return nil, err
	}
	model, err := t.issueToken(user, ClaimRefreshSub, ttl, nil)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (t *TokenService) VerifyToken(signedToken string) (*base.Token, error) {
	model, err := t.tokenRepository.FindTokenBySignedString(signedToken)
	if err != nil {
		return nil, err
	}

	token, err := t.jwt.Parse(signedToken)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims := token.Claims.(base.MapClaims)

	if claims["sub"] != model.Subject {
		return nil, errors.New("invalid token subject")
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, errors.New("token is expired")
	}

	return token, nil
}

func (t *TokenService) RefreshTokens(accessToken string, refreshToken string, options *TokenOptions) (*TokensResponse, error) {
	refreshModel, err := t.tokenRepository.FindTokenBySignedStringAndSubject(refreshToken, ClaimRefreshSub)
	if err != nil {
		return nil, err
	}

	_, err = t.VerifyToken(refreshToken)
	if err != nil {
		return nil, err
	}

	accessModel, err := t.tokenRepository.FindTokenBySignedStringAndSubject(accessToken, ClaimAccessSub)
	if err != nil {
		_ = t.RevokeToken(refreshToken)
		return nil, err
	}

	if *accessModel.RefreshTokenId != refreshModel.ID {
		return nil, errors.New("tokens pair does not math")
	}

	user := refreshModel.User

	_ = t.RevokeToken(accessToken)
	_ = t.RevokeToken(refreshToken)

	return t.IssueTokens(user, options)
}

func (t *TokenService) RevokeUserTokens(user *models.User) error {
	return t.tokenRepository.DeleteTokensByUID(user.UID)
}

func (t *TokenService) RevokeToken(signedToken string) error {
	model, err := t.tokenRepository.FindTokenBySignedString(signedToken)
	if err != nil {
		return err
	}

	if model.RefreshTokenId != nil {
		err = t.tokenRepository.DeleteTokenByID(*model.RefreshTokenId)
		if err != nil {
			return err
		}
	}

	return t.tokenRepository.Delete(model)
}

func (t *TokenService) issueToken(user *models.User, subject string, expire time.Duration, refreshId *uint64) (*models.Token, error) {
	exp := time.Now().Add(expire)

	claims := base.MapClaims{
		"sub":       subject,
		"exp":       exp.Unix(),
		"uid":       user.UID,
		"roleName":  user.RoleName,
		"parentId":  user.ParentId,
		"username":  user.Username,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
	}

	token := t.jwt.Issue(claims)
	jwtSigned, err := t.jwt.Sign(token)
	if err != nil {
		return nil, err
	}

	model := &models.Token{
		Subject:        subject,
		SignedString:   jwtSigned,
		UserUID:        user.UID,
		RefreshTokenId: refreshId,
	}

	created, err := t.tokenRepository.Create(model)
	if err != nil {
		return nil, err
	}

	return created, nil
}
