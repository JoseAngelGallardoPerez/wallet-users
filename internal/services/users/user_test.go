package users

import (
	"context"
	"database/sql/driver"
	"testing"

	pbSettings "github.com/Confialink/wallet-settings/rpc/proto/settings"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/services"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
	"github.com/Confialink/wallet-users/internal/services/syssettings/mocks"
	"github.com/Confialink/wallet-users/internal/tests/mocks/helpers"
	"github.com/Confialink/wallet-users/internal/tests/mocks/vendor-mocks/rpc/settings"
)

type AnyValue struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyValue) Match(v driver.Value) bool {
	return true
}

func TestUserCreateByRequest(t *testing.T) {
	gormMock := helpers.DbMock.GetGormMock()
	dbMock := helpers.DbMock.GetDbMock()

	userRepo := repositories.NewUsersRepository(gormMock)
	addressRepo := repositories.NewAddressRepository(gormMock)
	attributeRepo := repositories.NewAttributeRepository(gormMock)
	userAttributeValueRepo := repositories.NewUserAttributeValueRepository(gormMock)
	attributeService := NewAttributeService(attributeRepo, userAttributeValueRepo)

	client := &settings.MockSettingsHandler{}
	req := &pbSettings.Request{Path: "regional/modules/velmie_wallet_gdpr"}
	resp := &pbSettings.Response{Setting: &pbSettings.Setting{Value: "disabled"}}
	client.On("Get", context.Background(), req).Return(resp, nil)
	clientFactory := &mocks.ClientFactory{}
	clientFactory.On("NewClient").Return(client, nil)
	companyRepo := repositories.NewCompanyRepository(gormMock)
	companyService := NewCompanyService(companyRepo)

	service := NewUserService(
		gormMock,
		userRepo,
		companyService,
		addressRepo,
		attributeService,
		nil,
		services.NewPassword(),
		nil,
		nil,
		syssettings.NewSysSettings(clientFactory),
		nil,
		nil,
	)

	uid := "testUid"
	attributeName := "testAttr"
	attributeValue := "testValue"
	email := "example@example.com"
	user := &models.User{
		Email: email,
		PhysicalAddresses: []*models.Address{
			{
				Type: models.AddressTypePhysical,
				City: "test",
			},
		},
		MailingAddresses: []*models.Address{
			{
				Type: models.AddressTypeMailing,
				City: "test2",
			},
		},
		Attributes: map[string]interface{}{attributeName: attributeValue},
	}

	// Mock start transaction
	dbMock.ExpectBegin()

	// Mock insert User model query
	dbMock.ExpectExec("INSERT INTO `users`").
		WithArgs(AnyValue{}, email, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}).WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock select queries
	sqlRows := sqlmock.NewRows([]string{"uid"}).AddRow(uid)
	dbMock.ExpectQuery("^SELECT (.+) FROM `users` WHERE (.+)").WillReturnRows(sqlRows)
	sqlRows = sqlmock.NewRows([]string{"uid"}).AddRow(uid)
	dbMock.ExpectQuery("^SELECT (.+) FROM `users` WHERE (.+)").WillReturnRows(sqlRows)

	// Mock insert MailingAddresses
	dbMock.ExpectExec("INSERT INTO `addresses`").
		WithArgs(uid, models.AddressTypeMailing, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}).WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock insert PhysicalAddresses
	dbMock.ExpectExec("INSERT INTO `addresses`").
		WithArgs(uid, models.AddressTypePhysical, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}, AnyValue{}).WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock Attributes queries
	attributeId := 1
	sqlRows = sqlmock.NewRows([]string{"id", "slug"}).AddRow(attributeId, attributeName)
	dbMock.ExpectQuery("^SELECT (.+) FROM `attributes` WHERE (.+)").WillReturnRows(sqlRows)
	dbMock.ExpectExec("UPDATE `user_attribute_values`").
		WithArgs(attributeValue, AnyValue{}, uid, attributeId).WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock commit transaction
	dbMock.ExpectCommit()

	_, err := service.CreateNew(user)
	assert.Nil(t, err, "service must not return an error")
}
