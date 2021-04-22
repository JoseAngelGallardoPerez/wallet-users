package forms

import "github.com/pkg/errors"

const (
	fieldTypeArray         = "array"
	fieldTypeObject        = "object"
	fieldTypeString        = "string"
	fieldTypeStringPointer = "stringPointer"
	fieldTypeInt           = "int"
	fieldTypeIntPointer    = "intPointer"
	fieldTypeFloat         = "float"
	fieldTypeBool          = "bool"
)

type FieldType interface {
	// Returns default value for type
	defaultValue() interface{}
}

// Returns default value for particular type
func defaultValueForType(fieldType string) (FieldType, error) {

	switch fieldType {
	case fieldTypeString:
		return NewString(), nil
	case fieldTypeStringPointer:
		return NewStringPointer(), nil
	case fieldTypeInt:
		return NewInt(), nil
	case fieldTypeIntPointer:
		return NewIntPointer(), nil
	case fieldTypeFloat:
		return NewFloat(), nil
	case fieldTypeBool:
		return NewBool(), nil
	}

	return nil, errors.Errorf("invalid field type: %s", fieldType)
}

type String struct {
	value string
}

func NewString() FieldType {
	return &String{}
}

func (s *String) defaultValue() interface{} {
	return s.value
}

type StringPointer struct {
	value *string
}

func NewStringPointer() FieldType {
	return &StringPointer{}
}

func (s *StringPointer) defaultValue() interface{} {
	return s.value
}

type Int struct {
	value int
}

func NewInt() FieldType {
	return &Int{}
}

func (s *Int) defaultValue() interface{} {
	return s.value
}

type IntPointer struct {
	value *int
}

func NewIntPointer() FieldType {
	return &IntPointer{}
}

func (s *IntPointer) defaultValue() interface{} {
	return s.value
}

type Float struct {
	value float64
}

func NewFloat() FieldType {
	return &Float{}
}

func (s *Float) defaultValue() interface{} {
	return s.value
}

type Bool struct {
	value bool
}

func NewBool() FieldType {
	return &Bool{}
}

func (s *Bool) defaultValue() interface{} {
	return s.value
}
