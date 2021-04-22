package forms

import (
	"context"
	"fmt"
	"testing"

	pbSettings "github.com/Confialink/wallet-settings/rpc/proto/settings"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	validatorPkg "gopkg.in/go-playground/validator.v8"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	form_conditions "github.com/Confialink/wallet-users/internal/http/form-conditions"
	"github.com/Confialink/wallet-users/internal/http/responses"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
	"github.com/Confialink/wallet-users/internal/services/syssettings/mocks"
	"github.com/Confialink/wallet-users/internal/tests/mocks/helpers"
	"github.com/Confialink/wallet-users/internal/tests/mocks/log15"
	"github.com/Confialink/wallet-users/internal/tests/mocks/vendor-mocks/rpc/settings"
)

func TestUserSignUp(t *testing.T) {
	gormMock := helpers.DbMock.GetGormMock()
	dbMock := helpers.DbMock.GetDbMock()
	repo := repositories.NewFormRepository(gormMock)
	factory := NewFactory(repo)

	form := &models.Form{
		Type:           formTypeSignUp,
		OwnerRoleNames: `["buyer"]`,
		Form: `{
"conditions": [
    {
      "name": "useDefaultClass"
    }
  ],
  "fields": [
    {
      "name": "email",
      "type": "string",
      "validators": [
        {
          "name": "email",
          "param": ""
        },
        {
          "name": "max",
          "param": "255"
        }
      ]
    },
	{
      "name": "password",
      "type": "string",
      "validators": [
        {
          "name": "required"
        },
        {
          "name": "max",
          "param": "255"
        }
      ]
    },
	{
      "name": "roleName",
      "type": "string",
      "validators": [
        {
          "name": "max",
          "param": "255"
        }
      ]
    },
    {
      "name": "physicalAddresses",
      "type": "array",
      "validators": [
        {
          "name": "required",
          "param": ""
        },
        {
          "name": "max",
          "param": "2"
        },
        {
          "name": "dive",
          "param": ""
        }
      ],
      "children": [
        {
          "name": "city",
          "type": "string",
          "validators": [
            {
              "name": "required",
              "param": ""
            },
            {
              "name": "max",
              "param": "255"
            }
          ]
        }
      ]
    },
    {
      "name": "attributes",
      "type": "object",
      "validators": [
        {
          "name": "required",
          "param": ""
        },
        {
          "name": "max",
          "param": "2"
        },
        {
          "name": "dive",
          "param": ""
        }
      ],
      "children": [
        {
          "name": "prop1",
          "type": "string",
          "validators": [
            {
              "name": "required",
              "param": ""
            },
            {
              "name": "max",
              "param": "255"
            }
          ]
        }
      ]
    }
  ]
}`,
	}

	sqlRows := sqlmock.NewRows([]string{"type", "owner_role_names", "form"}).
		AddRow(form.Type, form.OwnerRoleNames, form.Form)
	dbMock.ExpectQuery("^SELECT (.+) FROM `forms`$").WillReturnRows(sqlRows)
	err := factory.InitForms()
	assert.Nil(t, err, "factory must not return an error")

	client := &settings.MockSettingsHandler{}
	roleName := "buyer"
	req := &pbSettings.Request{Path: "profile/default-user-classes/" + roleName}
	defaultUserClassId := "2"
	resp := &pbSettings.Response{Setting: &pbSettings.Setting{Value: defaultUserClassId}}
	client.On("Get", context.Background(), req).Return(resp, nil)
	clientFactory := &mocks.ClientFactory{}
	clientFactory.On("NewClient").Return(client, nil)
	defaultUserClassService := form_conditions.NewDefaultUserClass(syssettings.NewSysSettings(clientFactory))
	conditionRegistry := form_conditions.NewConditionRegistry()
	err = conditionRegistry.Register(defaultUserClassService)
	assert.Nil(t, err, "service must not return an error")
	logger := logger.Logger{}
	service := NewUser(
		validatorPkg.New(&validatorPkg.Config{TagName: "binding"}),
		syssettings.NewSysSettings(clientFactory),
		factory,
		conditionRegistry,
		&logger,
	)

	/*
		Invalid input data
	*/
	rawData := `{"roleName": "buyer"` // invalid json
	_, err = service.SignUp([]byte(rawData))
	assert.Error(t, err, "service must return an error")
	errorText := "cannot unmarshal body"
	assert.Contains(t, err.Error(), errorText, fmt.Sprintf("error must contain the text: `%s`", errorText))

	/*
		Empty role
	*/
	rawData = `{"roleName2": "buyer"}`
	_, err = service.SignUp([]byte(rawData))
	assert.IsType(t, validatorPkg.ValidationErrors{}, err, "Service must return validation errors")

	/*
		Invalid role. Form is not exist
	*/
	rawData = `{"roleName": "random_role"}`
	source := &logger
	source.On("Error", "cannot make form", "role", "random_role", "formId", "sign_up_random_role_", "err", mock.Anything).Return()
	errorText = responses.UnsupportedRole
	_, err = service.SignUp([]byte(rawData))
	assert.Error(t, err, "service must return an error")
	assert.Contains(t, err.Error(), errorText, fmt.Sprintf("error must contain the text: `%s`", errorText))

	/*
		Successful case
	*/
	email := "test@example.com"
	city1 := "city_1"
	prop1 := "value"
	password := "pass"

	rawData = fmt.Sprintf(`{"roleName": "%s","password": "%s","email":"%s","physicalAddresses":[{"city":"%s"}],"attributes":{"prop1":"%s"}}`, roleName, password, email, city1, prop1)

	user, err := service.SignUp([]byte(rawData))
	assert.Nil(t, err, "service must not return an error")
	assert.Equal(t, user.Email, email, "Emails does not match")
	assert.Equal(t, user.Password, password, "Password does not match")
	assert.Equal(t, user.RoleName, roleName, "Roles does not match")
	assert.Equal(t, user.Attributes["prop1"], prop1, "Attributes.Prop1 does not match")
	assert.Equal(t, user.PhysicalAddresses[0].City, city1, "PhysicalAddresses[0].City does not match")
	assert.Equal(t, user.Status, models.StatusActive, "Statuses does not match")
	assert.Equal(t, string(user.ClassId), defaultUserClassId, "userClass does not match")

	/*
		Password field not found
	*/
	rawData = fmt.Sprintf(`{"roleName": "%s","email":"%s","physicalAddresses":[{"city":"%s"}],"attributes":{"prop1":"%s"}}`, roleName, email, city1, prop1)
	_, err = service.SignUp([]byte(rawData))
	assert.Error(t, err, "service must return an error")
}
