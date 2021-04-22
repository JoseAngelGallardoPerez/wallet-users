package forms

import (
	"encoding/json"
	"fmt"
	"strings"

	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"github.com/pkg/errors"

	"github.com/Confialink/wallet-users/internal/db/repositories"
)

const (
	formTypeSignUp = "sign_up"
	formTypeUpdate = "update"
)

// This configuration describes fields of a form and additional conditions
type FormConfig struct {
	// List of fields
	Fields []*Field `json:"fields"`

	// Additional conditions for implementation logic
	Conditions []*Condition `json:"conditions"`
}

type Condition struct {
	Name   string      `json:"name"`
	Params interface{} `json:"params"`
}

// Field of a form
type Field struct {
	Name       string       `json:"name"`
	Type       string       `json:"type"`
	Validators []*Validator `json:"validators"`
	Children   []*Field     `json:"children"`
}

// Validator of a field
type Validator struct {
	Name  string `json:"name"`
	Param string `json:"param"`
}

type Form struct {
	// It build a new structure to validate input data
	FormBuilder dynamicstruct.DynamicStruct

	// Additional conditions for implementation logic
	Conditions []*Condition
}

type Factory struct {
	repo  *repositories.FormRepository
	cache map[string]*Form
}

// Constructor
func NewFactory(repo *repositories.FormRepository) *Factory {
	return &Factory{cache: make(map[string]*Form), repo: repo}
}

// Obtains a form builder from cache
func (s *Factory) Form(formId string) (*Form, error) {
	form, ok := s.cache[formId]
	if ok {
		return form, nil
	}

	return nil, errors.Errorf("cannot find form by formId: %s", formId)
}

// Initializes all available forms.
// Prepares form builders and puts them into the cache.
func (s *Factory) InitForms() error {
	forms, err := s.repo.All()
	if err != nil {
		return errors.Wrap(err, "cannot receive forms from DB")
	}

	for _, form := range forms {
		ownerRoles, err := form.OwnerRoleNamesAsList()
		if err != nil {
			return errors.Wrap(err, "cannot receive list of ownerRoles")
		}

		initiatorRoles, err := form.InitiatorRoleNamesAsList()
		if err != nil {
			return errors.Wrap(err, "cannot receive list of initiatorRoles")
		}

		if len(ownerRoles) == 0 {
			return errors.Errorf("OwnerRoleNames is empty. formId: %d", form.ID)
		}

		// Register a form builder in the cache by unique Id.
		// Form Id is a mix of role names
		// formId = {formType}_{ownerRole}_{initiatorRole}
		// Example: `sign_up_buyer_admin` - When an admin tries to update an entity which belongs to buyer
		// Example: `sign_up_buyer_` - When an anonymous user tries to update an entity which belongs to buyer
		for _, ownerRole := range ownerRoles {
			if len(form.Type) == 0 {
				return errors.Errorf("Type is empty. formId: %d", form.ID)
			}
			formId := form.Type + "_" + ownerRole
			if len(initiatorRoles) > 0 {
				for _, initiatorRole := range initiatorRoles {
					innerFormId := formId + "_" + initiatorRole
					if err := s.prepareForm(form.Form, innerFormId); err != nil {
						return err
					}
				}
			} else {
				formId += "_"
				if err := s.prepareForm(form.Form, formId); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Prepares a form builder and put it into the cache
func (s *Factory) prepareForm(jsonConfiguration, formId string) error {
	formConfig := &FormConfig{}

	// Fills the form by raw data
	if err := json.Unmarshal([]byte(jsonConfiguration), formConfig); err != nil {
		return errors.Wrapf(err, "form `%s`.cannot unmarshal configuration: %s", formId, jsonConfiguration)
	}

	if len(formConfig.Fields) == 0 {
		return errors.Errorf("the form `%s` does not contain fields: %s", formId, jsonConfiguration)
	}

	formBuilder, err := s.buildForm(formConfig.Fields)
	if err != nil {
		return errors.Wrapf(err, "form `%s`", formId)
	}

	s.cache[formId] = &Form{formBuilder, formConfig.Conditions}

	return nil
}

// Returns a form builder by list of fields
func (s *Factory) buildForm(fields []*Field) (dynamicstruct.DynamicStruct, error) {
	parentStruct := dynamicstruct.NewStruct()
	for _, field := range fields {
		switch field.Type {
		case fieldTypeArray:
			// Array field.
			// Contains list of structures(objects).
			if len(field.Children) == 0 {
				return nil, errors.Errorf("children are empty. field: %s", field.Name)
			}
			arrayStruct, err := s.buildForm(field.Children)
			if err != nil {
				return nil, err
			}
			definition := dynamicstruct.ExtendStruct(arrayStruct.New()).Build()
			childStruct := definition.NewSliceOfStructs()
			s.addField(parentStruct, field, childStruct)
		case fieldTypeObject:
			// Object field.
			// Contains a structure
			if len(field.Children) == 0 {
				return nil, errors.Errorf("children are empty. field: %s", field.Name)
			}
			objectStruct, err := s.buildForm(field.Children)
			if err != nil {
				return nil, err
			}
			s.addField(parentStruct, field, objectStruct.New())
		default:
			// Scalar field
			fieldType, err := defaultValueForType(field.Type)
			if err != nil {
				return nil, errors.Wrapf(err, "field: %s", field.Name)
			}

			s.addField(parentStruct, field, fieldType.defaultValue())
		}
	}

	return parentStruct.Build(), nil
}

// Add a new filed to the builder
func (s *Factory) addField(builder dynamicstruct.Builder, field *Field, fieldType interface{}) {
	tags := fmt.Sprintf(`json:"%s"`, field.Name)
	tags += fmt.Sprintf(` binding:"%s"`, s.validatorsAsString(field.Validators))

	builder.AddField(
		strings.Title(field.Name),
		fieldType,
		tags,
	)
}

// Returns list of validators as string.
// Example `required,max=2,dive`
func (s *Factory) validatorsAsString(vldrs []*Validator) string {
	validators := make([]string, 0, len(vldrs))
	for _, validator := range vldrs {
		validatorString := validator.Name
		if len(validator.Param) > 0 {
			validatorString += "=" + validator.Param
		}
		validators = append(validators, validatorString)
	}

	return strings.Join(validators, ",")
}
