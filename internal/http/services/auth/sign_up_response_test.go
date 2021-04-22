package auth

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dgrijalva/jwt-go"
	base "github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/Confialink/wallet-users/internal/config"
	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	jwtMocks "github.com/Confialink/wallet-users/internal/jwt/mocks"
	"github.com/Confialink/wallet-users/internal/services/auth"
	"github.com/Confialink/wallet-users/internal/services/auth/mocks"
	"github.com/Confialink/wallet-users/internal/tests/mocks/helpers"
)

type AnyValue struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyValue) Match(v driver.Value) bool {
	return true
}

var _ = Describe("services package", func() {
	gormMock := helpers.DbMock.GetGormMock()
	dbMock := helpers.DbMock.GetDbMock()
	config := &config.Configuration{JWT: &config.JwtConfiguration{SigningMethod: jwt.SigningMethodHS512, Secret: "secret"}}
	tokenRepository := repositories.NewTokenRepository(gormMock)
	tmpTokensService := auth.NewTemporaryTokens(&jwtMocks.Service{})

	Context("Make", func() {
		Context("a user is active", func() {
			user := &models.User{Status: models.StatusActive}
			When("tokenService cannot issues tokens", func() {
				It("returns an error", func() {
					ttlResolverMock := &mocks.TokenTTLResolver{}
					ttlResolverMock.
						On("ResolveByTokenSubject", auth.ClaimRefreshSub).
						Return(time.Duration(1*time.Hour), errors.New("error text"))
					tokenService := auth.TokenServiceFactory(config, tokenRepository, ttlResolverMock, nil)
					service := NewSignUpResponse(tokenService, tmpTokensService)

					res, err := service.Make(user)
					Expect(res).Should(BeNil())
					Expect(err).Should(HaveOccurred())
				})
			})

			When("everything is OK", func() {
				It("returns a response", func() {
					ttlResolverMock := &mocks.TokenTTLResolver{}
					ttlResolverMock.
						On("ResolveByTokenSubject", auth.ClaimRefreshSub).
						Return(time.Duration(1*time.Hour), nil)
					ttlResolverMock.
						On("ResolveByTokenSubject", auth.ClaimAccessSub).
						Return(time.Duration(1*time.Hour), nil)
					tokenService := auth.TokenServiceFactory(config, tokenRepository, ttlResolverMock, nil)
					service := NewSignUpResponse(tokenService, tmpTokensService)

					// mock inserting refresh token
					dbMock.ExpectBegin()
					dbMock.ExpectExec("^INSERT INTO `tokens`").
						WithArgs(AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}).
						WillReturnResult(sqlmock.NewResult(1, 1))
					dbMock.ExpectCommit()

					// mock inserting access token
					dbMock.ExpectBegin()
					dbMock.ExpectExec("^INSERT INTO `tokens`").
						WithArgs(AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}).
						WillReturnResult(sqlmock.NewResult(1, 1))
					dbMock.ExpectCommit()

					res, err := service.Make(user)
					Expect(len(res.Headers)).Should(Equal(0))

					_, ok := res.Data.(*auth.TokensResponse)
					Expect(ok).Should(BeTrue())
					Expect(err).ShouldNot(HaveOccurred())
					Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
				})
			})
		})

		Context("Make", func() {
			Context("a user is not active", func() {
				user := &models.User{Status: models.StatusPending}
				When("tmpTokensService cannot issue a token", func() {
					It("returns an error", func() {
						ttlResolverMock := &mocks.TokenTTLResolver{}
						jwtMock := &jwtMocks.Service{}
						tmpTokensService := auth.NewTemporaryTokens(jwtMock)
						baseToken := &base.Token{}
						jwtMock.On("Issue", mock.Anything).Return(baseToken)
						jwtMock.On("Sign", baseToken).Return("", errors.New("text"))

						tokenService := auth.TokenServiceFactory(config, tokenRepository, ttlResolverMock, nil)
						service := NewSignUpResponse(tokenService, tmpTokensService)

						res, err := service.Make(user)
						Expect(res).Should(BeNil())
						Expect(err).Should(HaveOccurred())
					})
				})
				When("everything is OK", func() {
					It("returns a response", func() {
						ttlResolverMock := &mocks.TokenTTLResolver{}
						jwtMock := &jwtMocks.Service{}
						tmpTokensService := auth.NewTemporaryTokens(jwtMock)
						baseToken := &base.Token{}
						jwtMock.On("Issue", mock.Anything).Return(baseToken)
						jwtMock.On("Sign", baseToken).Return("text", nil)

						tokenService := auth.TokenServiceFactory(config, tokenRepository, ttlResolverMock, nil)
						service := NewSignUpResponse(tokenService, tmpTokensService)

						res, err := service.Make(user)
						Expect(len(res.Headers)).Should(Equal(2))

						_, ok := res.Data.(*models.User)
						Expect(ok).Should(BeTrue())
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})
		})
	})
})
