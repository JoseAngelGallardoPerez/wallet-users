package users

import (
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/tests/mocks/helpers"
)

var _ = Describe("users package", func() {

	uid := "test-uid"
	user := &models.User{UID: uid}
	gormMock := helpers.DbMock.GetGormMock()
	dbMock := helpers.DbMock.GetDbMock()
	repo := repositories.NewConfirmationCodeRepository(gormMock)
	service := NewConfirmationCode(repo)

	Context("CreateNewVerificationCode", func() {
		deleteQuery := "^DELETE FROM `confirmation_codes` WHERE \\(user_uid = (.+) AND subject = (.+)\\)"

		When("there is a some DB error during delete old confirmation code", func() {
			It("returns an error", func() {
				errorText := "cannot delete old code"
				dbMock.ExpectBegin()
				dbMock.ExpectExec(deleteQuery).
					WithArgs(uid, models.ConfirmationCodeSubjectPhoneVerificationCode).
					WillReturnError(errors.New(errorText))

				_, err := service.CreateNewVerificationCode(user, models.ConfirmationCodeSubjectPhoneVerificationCode)
				Expect(err).Should(MatchError(errorText))
				Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
			})
		})

		When("there is a some DB error during insert confirmation code", func() {
			It("returns an error", func() {
				errorText := "cannot insert new code"
				dbMock.ExpectBegin()
				dbMock.ExpectExec(deleteQuery).
					WithArgs(uid, models.ConfirmationCodeSubjectPhoneVerificationCode).
					WillReturnResult(sqlmock.NewResult(0, 0))
				dbMock.ExpectCommit()

				dbMock.ExpectBegin()
				dbMock.ExpectExec("^INSERT INTO `confirmation_codes`").
					WithArgs(AnyValue{}, uid, models.ConfirmationCodeSubjectPhoneVerificationCode, AnyValue{}, AnyValue{}, AnyValue{}).
					WillReturnError(errors.New(errorText))

				_, err := service.CreateNewVerificationCode(user, models.ConfirmationCodeSubjectPhoneVerificationCode)
				Expect(err).Should(MatchError(errorText))
				Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
			})
		})

		When("confirmation code is successfully generated", func() {
			It("does not return an error", func() {
				dbMock.ExpectBegin()
				dbMock.ExpectExec(deleteQuery).
					WithArgs(uid, models.ConfirmationCodeSubjectPhoneVerificationCode).
					WillReturnResult(sqlmock.NewResult(0, 0))
				dbMock.ExpectCommit()

				dbMock.ExpectBegin()
				dbMock.ExpectExec("^INSERT INTO `confirmation_codes`").
					WithArgs(AnyValue{}, uid, models.ConfirmationCodeSubjectPhoneVerificationCode, AnyValue{}, AnyValue{}, AnyValue{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				dbMock.ExpectCommit()

				code, err := service.CreateNewVerificationCode(user, models.ConfirmationCodeSubjectPhoneVerificationCode)
				Expect(err).Should(BeNil())
				Expect(dbMock.ExpectationsWereMet()).Should(BeNil())

				Expect(code.UserUID).Should(Equal(uid))
				Expect(code.Subject).Should(Equal(models.ConfirmationCodeSubjectPhoneVerificationCode))
			})
		})
	})

	code := "random-code"
	selectQuery := "SELECT (.+) FROM `confirmation_codes` WHERE \\(user_uid = \\?\\) AND \\(subject = \\?\\) AND \\(code = \\?\\) AND \\(expires_at >= \\?\\) (.+)"

	Context("CheckPhoneCode", func() {
		When("there is a some DB error during select query", func() {
			It("returns an error", func() {
				errorText := "db error"
				dbMock.ExpectQuery(selectQuery).
					WithArgs(uid, models.ConfirmationCodeSubjectPhoneVerificationCode, code, AnyValue{}).
					WillReturnError(errors.New(errorText))

				Expect(service.CheckPhoneCode(code, user)).Should(MatchError(errorText))
				Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
			})
		})

		When("code is successfully checked", func() {
			It("returns nil", func() {
				sqlRows := sqlmock.NewRows([]string{"id"}).AddRow("1")
				dbMock.ExpectQuery(selectQuery).
					WithArgs(uid, models.ConfirmationCodeSubjectPhoneVerificationCode, code, AnyValue{}).
					WillReturnRows(sqlRows)

				Expect(service.CheckPhoneCode(code, user)).Should(BeNil())
				Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
			})
		})
	})

	Context("CheckEmailCode", func() {
		When("there is a some DB error during select query", func() {
			It("returns an error", func() {
				errorText := "db error"
				dbMock.ExpectQuery(selectQuery).
					WithArgs(uid, models.ConfirmationCodeSubjectEmailVerificationCode, code, AnyValue{}).
					WillReturnError(errors.New(errorText))

				Expect(service.CheckEmailCode(code, user)).Should(MatchError(errorText))
				Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
			})
		})

		When("code is successfully checked", func() {
			It("returns nil", func() {
				sqlRows := sqlmock.NewRows([]string{"id"}).AddRow("1")
				dbMock.ExpectQuery(selectQuery).
					WithArgs(uid, models.ConfirmationCodeSubjectEmailVerificationCode, code, AnyValue{}).
					WillReturnRows(sqlRows)

				Expect(service.CheckEmailCode(code, user)).Should(BeNil())
				Expect(dbMock.ExpectationsWereMet()).Should(BeNil())
			})
		})
	})
})
