package forms

import (
	"encoding/json"
	"fmt"
	"time"

	perrors "github.com/Confialink/wallet-pkg-errors"
	"github.com/inconshreveable/log15"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"github.com/pkg/errors"

	userpb "github.com/Confialink/wallet-users/rpc/proto/users"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/http/form-conditions"
	"github.com/Confialink/wallet-users/internal/http/responses"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
	"github.com/Confialink/wallet-users/internal/validators"
)

type Role struct {
	RoleName string `json:"roleName" binding:"required"`
}

type User struct {
	validator         validators.Interface
	sysSettings       *syssettings.SysSettings
	formFactory       *Factory
	conditionRegistry *form_conditions.ConditionRegistry
	logger            log15.Logger
}

func NewUser(
	validator validators.Interface,
	sysSettings *syssettings.SysSettings,
	formFactory *Factory,
	conditionRegistry *form_conditions.ConditionRegistry,
	logger log15.Logger,
) *User {
	return &User{validator, sysSettings, formFactory, conditionRegistry, logger}
}

// Returns User model after the signup process.
// Validates and fills the User structure.
func (u *User) Update(user *models.User, editor *userpb.User, rawData []byte) error {

	formId := formTypeUpdate + "_" + user.RoleName + "_" + editor.RoleName

	form, err := u.formFactory.Form(formId)
	if err != nil {
		return errors.Wrapf(err, "cannot make form. role: %s, formId: %s", user.RoleName, formId)
	}
	structure := form.FormBuilder.New()

	// Fill the form by raw data
	if err := json.Unmarshal(rawData, structure); err != nil {
		return errors.Wrapf(err, "cannot unmarshal data. data: %s", rawData)
	}

	// Add UID of current user to the struct. It is required for validators
	uidRawData := []byte(fmt.Sprintf(`{"uid": "%s"}`, user.UID))
	if err := json.Unmarshal(uidRawData, structure); err != nil {
		return errors.Wrapf(err, "cannot unmarshal data. cannot add uid. data: %s", rawData)
	}

	// We wrap our dynamic structure to create a namespace.
	// Validator uses this namespace to detect 'source' of error field
	validatedStruct := struct {
		InnerStruct interface{} `binding:"dive"`
	}{structure}

	// Validate the structure
	if err := u.validator.Struct(validatedStruct); err != nil {
		return err
	}

	// Convert validated struct to a json string
	data, err := json.Marshal(structure)
	if err != nil {
		return errors.Wrapf(err, "cannot marshal. struct: %+v", structure)
	}

	// Fill the User model by validated data
	if err := json.Unmarshal(data, user); err != nil {
		return errors.Wrapf(err, "cannot unmarshal form data. struct: %s", data)
	}

	conditions, err := u.conditions(formId)
	if err != nil {
		return err
	}

	if len(conditions) > 0 {
		for _, condition := range conditions {
			conditionService, err := u.conditionRegistry.Make(condition.Name)
			if err != nil {
				return err
			}
			if err := conditionService.Apply(user, condition.Params); err != nil {
				return err
			}
		}
	}

	return nil
}

// Returns User model after the signup process.
// Validates and fills the User structure.
func (u *User) SignUp(rawData []byte) (*models.User, error) {
	roleForm := &Role{}
	if err := json.Unmarshal(rawData, roleForm); err != nil {
		return nil, errors.Wrapf(err, "cannot unmarshal body: %s", rawData)
	}
	if err := u.validator.Struct(roleForm); err != nil {
		return nil, err
	}

	formId := formTypeSignUp + "_" + roleForm.RoleName + "_"

	user := &models.User{
		RoleName:         roleForm.RoleName,
		Status:           models.StatusActive,
		IsEmailConfirmed: false,
		IsPhoneConfirmed: false,
	}
	user.LastActedAt = time.Now()
	user, err := u.update(user, rawData, formId)
	if err != nil {
		return nil, err
	}

	conditions, err := u.conditions(formId)
	if err != nil {
		return nil, err
	}

	if len(conditions) > 0 {
		for _, condition := range conditions {
			conditionService, err := u.conditionRegistry.Make(condition.Name)
			if err != nil {
				return nil, err
			}
			if err := conditionService.Apply(user, condition.Params); err != nil {
				return nil, err
			}
		}
	}

	return user, nil
}

// Returns additional form conditions
func (u *User) conditions(formId string) ([]*Condition, error) {
	form, err := u.formFactory.Form(formId)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot make form. formId: %s", formId)
	}

	return form.Conditions, nil
}

// Updates user fields in the structure.
// Result depends on roles and registered forms.
func (u *User) update(user *models.User, rawData []byte, formId string) (*models.User, error) {
	form, err := u.formFactory.Form(formId)
	if err != nil {
		u.logger.Error("cannot make form", "role", user.RoleName, "formId", formId, "err", err)
		return nil, &perrors.PublicError{Title: "Unknown role name", Code: responses.UnsupportedRole, HttpStatus: 400, Meta: ""}
	}
	structure := form.FormBuilder.New()

	// Fill the form by raw data
	if err := json.Unmarshal(rawData, structure); err != nil {
		return nil, errors.Wrapf(err, "cannot unmarshal data. data: %s", rawData)
	}

	// We wrap our dynamic structure to create a namespace.
	// Validator uses this namespace to detect 'source' of error field
	validatedStruct := struct {
		InnerStruct interface{} `binding:"dive"`
	}{structure}

	// Validate the structure
	if err := u.validator.Struct(validatedStruct); err != nil {
		return nil, err
	}

	// Convert validated struct to a json string
	data, err := json.Marshal(structure)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot marshal. struct: %+v", structure)
	}

	// Fill the User model by validated data
	if err := json.Unmarshal(data, user); err != nil {
		return nil, errors.Wrapf(err, "cannot unmarshal form data. struct: %s", data)
	}

	// Password field is not exported in the Users model.
	// We should read and set that manually
	field := dynamicstruct.NewReader(structure).GetField("Password")
	if field == nil {
		return nil, errors.New("Password field not found in structure")
	}

	user.Password = field.String()

	return user, nil
}
