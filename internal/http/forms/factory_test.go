package forms

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/tests/mocks/helpers"
)

func TestFactoryInitForms(t *testing.T) {
	gormMock := helpers.DbMock.GetGormMock()
	dbMock := helpers.DbMock.GetDbMock()
	repo := repositories.NewFormRepository(gormMock)
	service := NewFactory(repo)

	/*
		DB error
	*/
	err := service.InitForms()
	assert.Error(t, err, "service must return an error")
	errorText := "cannot receive forms from DB"
	assert.Contains(t, err.Error(), errorText, fmt.Sprintf("error must contain the text: `%s`", errorText))

	/*
		Bad OwnerRoleNames list
	*/
	form := &models.Form{
		OwnerRoleNames: `["buyer"`, // syntax error in json
	}
	sqlRows := sqlmock.NewRows([]string{"owner_role_names"}).
		AddRow(form.OwnerRoleNames)
	dbMock.ExpectQuery("^SELECT (.+) FROM `forms`$").WillReturnRows(sqlRows)
	err = service.InitForms()
	assert.Error(t, err, "service must return an error")
	errorText = "cannot receive list of ownerRoles"
	assert.Contains(t, err.Error(), errorText, fmt.Sprintf("error must contain the text: `%s`", errorText))

	/*
		Bad InitiatorRoleNames list
	*/
	form = &models.Form{
		OwnerRoleNames:     `["admin"]`,
		InitiatorRoleNames: `["buyer"`, // syntax error in json
	}
	sqlRows = sqlmock.NewRows([]string{"initiator_role_names", "owner_role_names"}).
		AddRow(form.InitiatorRoleNames, form.OwnerRoleNames)
	dbMock.ExpectQuery("^SELECT (.+) FROM `forms`$").WillReturnRows(sqlRows)
	err = service.InitForms()
	assert.Error(t, err, "service must return an error")
	errorText = "cannot receive list of initiatorRoles"
	assert.Contains(t, err.Error(), errorText, fmt.Sprintf("error must contain the text: `%s`", errorText))

	/*
		OwnerRoleNames is empty
	*/
	form = &models.Form{InitiatorRoleNames: `["buyer"]`}
	sqlRows = sqlmock.NewRows([]string{"initiator_role_names"}).
		AddRow(form.InitiatorRoleNames)
	dbMock.ExpectQuery("^SELECT (.+) FROM `forms`$").WillReturnRows(sqlRows)
	err = service.InitForms()
	assert.Error(t, err, "service must return an error")
	errorText = "OwnerRoleNames is empty"
	assert.Contains(t, err.Error(), errorText, fmt.Sprintf("error must contain the text: `%s`", errorText))

	/*
		FormType is empty
	*/
	form = &models.Form{OwnerRoleNames: `["buyer"]`}
	sqlRows = sqlmock.NewRows([]string{"owner_role_names"}).
		AddRow(form.OwnerRoleNames)
	dbMock.ExpectQuery("^SELECT (.+) FROM `forms`$").WillReturnRows(sqlRows)
	err = service.InitForms()
	assert.Error(t, err, "service must return an error")
	errorText = "Type is empty"
	assert.Contains(t, err.Error(), errorText, fmt.Sprintf("error must contain the text: `%s`", errorText))

	/*
		Invalid form configuration
	*/
	form = &models.Form{
		Type:           formTypeSignUp,
		OwnerRoleNames: `["buyer"]`,
		Form:           `{"fields": [{"name": "email","type": "string"}]`, // syntax error in json
	}
	sqlRows = sqlmock.NewRows([]string{"type", "owner_role_names", "form"}).
		AddRow(form.Type, form.OwnerRoleNames, form.Form)
	dbMock.ExpectQuery("^SELECT (.+) FROM `forms`$").WillReturnRows(sqlRows)
	err = service.InitForms()
	assert.Error(t, err, "service must return an error")
	errorText = "cannot unmarshal configuration"
	assert.Contains(t, err.Error(), errorText, fmt.Sprintf("error must contain the text: `%s`", errorText))

	/*
		The list of fields is empty
	*/
	form = &models.Form{
		Type:           formTypeSignUp,
		OwnerRoleNames: `["buyer"]`,
		Form:           `{"fields": []}`,
	}
	sqlRows = sqlmock.NewRows([]string{"type", "owner_role_names", "form"}).
		AddRow(form.Type, form.OwnerRoleNames, form.Form)
	dbMock.ExpectQuery("^SELECT (.+) FROM `forms`$").WillReturnRows(sqlRows)
	err = service.InitForms()
	assert.Error(t, err, "service must return an error")
	errorText = "does not contain fields"
	assert.Contains(t, err.Error(), errorText, fmt.Sprintf("error must contain the text: `%s`", errorText))

	/*
		An "array" field does not have any children fields
	*/
	form = &models.Form{
		Type:           formTypeSignUp,
		OwnerRoleNames: `["buyer"]`,
		Form:           `{"fields": [{"name": "email","type": "array"}]}`,
	}
	sqlRows = sqlmock.NewRows([]string{"type", "owner_role_names", "form"}).
		AddRow(form.Type, form.OwnerRoleNames, form.Form)
	dbMock.ExpectQuery("^SELECT (.+) FROM `forms`$").WillReturnRows(sqlRows)
	err = service.InitForms()
	assert.Error(t, err, "service must return an error")
	errorText = "children are empty"
	assert.Contains(t, err.Error(), errorText, fmt.Sprintf("error must contain the text: `%s`", errorText))

	/*
		An "object" field does not have any children fields
	*/
	form = &models.Form{
		Type:           formTypeSignUp,
		OwnerRoleNames: `["buyer"]`,
		Form:           `{"fields": [{"name": "email","type": "array"}]}`,
	}
	sqlRows = sqlmock.NewRows([]string{"type", "owner_role_names", "form"}).
		AddRow(form.Type, form.OwnerRoleNames, form.Form)
	dbMock.ExpectQuery("^SELECT (.+) FROM `forms`$").WillReturnRows(sqlRows)
	err = service.InitForms()
	assert.Error(t, err, "service must return an error")
	errorText = "children are empty"
	assert.Contains(t, err.Error(), errorText, fmt.Sprintf("error must contain the text: `%s`", errorText))

	/*
		A field has an unknown type
	*/
	form = &models.Form{
		Type:           formTypeSignUp,
		OwnerRoleNames: `["buyer"]`,
		Form:           `{"fields": [{"name": "email","type": "invalid_type"}]}`,
	}
	sqlRows = sqlmock.NewRows([]string{"type", "owner_role_names", "form"}).
		AddRow(form.Type, form.OwnerRoleNames, form.Form)
	dbMock.ExpectQuery("^SELECT (.+) FROM `forms`$").WillReturnRows(sqlRows)
	err = service.InitForms()
	assert.Error(t, err, "service must return an error")
	errorText = "invalid field type"
	assert.Contains(t, err.Error(), errorText, fmt.Sprintf("error must contain the text: `%s`", errorText))
	service.cache = make(map[string]*Form) // clear cache

	/*
		An successful initialization
	*/
	form = &models.Form{
		Type:           formTypeSignUp,
		OwnerRoleNames: `["buyer"]`,
		Form:           `{"fields": [{"name": "email","type": "string"}]}`,
	}
	sqlRows = sqlmock.NewRows([]string{"type", "owner_role_names", "form"}).
		AddRow(form.Type, form.OwnerRoleNames, form.Form)
	dbMock.ExpectQuery("^SELECT (.+) FROM `forms`$").WillReturnRows(sqlRows)
	err = service.InitForms()
	assert.Nil(t, err, "service must not return an error")
	assert.True(t, len(service.cache) == 1, "service must have one initialized form builder")
}

func TestFactoryForm(t *testing.T) {
	gormMock := helpers.DbMock.GetGormMock()
	dbMock := helpers.DbMock.GetDbMock()
	repo := repositories.NewFormRepository(gormMock)
	service := NewFactory(repo)

	/*
		Cannot find form
	*/
	_, err := service.Form("random-text")
	assert.Error(t, err, "service must return an error")
	errorText := "cannot find form"
	assert.Contains(t, err.Error(), errorText, fmt.Sprintf("error must contain the text: `%s`", errorText))

	/*
		A successful case
	*/
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
	err = service.InitForms()
	assert.Nil(t, err, "service must not return an error")

	builder, err := service.Form(formTypeSignUp + "_buyer_")
	assert.Nil(t, err, "Form method must not return an error")
	structure := builder.FormBuilder.New()

	email := "test@example.com"
	city1 := "city_1"
	prop1 := "value"

	rawData := fmt.Sprintf(`{"email":"%s","physicalAddresses":[{"city":"%s"}],"attributes":{"prop1":"%s"}}`, email, city1, prop1)

	// Fill the form by raw data
	if err := json.Unmarshal([]byte(rawData), structure); err != nil {
		assert.Nil(t, err, "invalid input data")
	}

	// Convert validated struct to a json string
	data, err := json.Marshal(structure)
	if err != nil {
		assert.Nil(t, err, "cannot Marshal the structure")
	}

	dataStruct1 := struct {
		Email      string `json:"email"`
		Attributes struct {
			Prop1 string `json:"prop1"`
		} `json:"attributes"`
		PhysicalAddresses []struct {
			City string `json:"city"`
		} `json:"physicalAddresses"`
	}{}
	dataStruct2 := struct {
		Email      string `json:"email"`
		Attributes struct {
			Prop1 string `json:"prop1"`
		} `json:"attributes"`
		PhysicalAddresses []struct {
			City string `json:"city"`
		} `json:"physicalAddresses"`
	}{}

	if err := json.Unmarshal([]byte(data), &dataStruct1); err != nil {
		assert.Nil(t, err, "invalid struct1 data")
	}
	if err := json.Unmarshal([]byte(rawData), &dataStruct2); err != nil {
		assert.Nil(t, err, "invalid struct2 data")
	}

	assert.Equal(t, dataStruct1.Email, dataStruct2.Email, "Emails does not match")
	assert.Equal(t, dataStruct1.Attributes.Prop1, dataStruct2.Attributes.Prop1, "Attributes.Prop1 does not match")
	assert.Equal(t, dataStruct1.PhysicalAddresses[0].City, dataStruct2.PhysicalAddresses[0].City, "PhysicalAddresses[0].City does not match")
}
