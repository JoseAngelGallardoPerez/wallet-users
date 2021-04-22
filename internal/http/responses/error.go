package responses

import (
	"net/http"

	"github.com/Confialink/wallet-pkg-errors"
)

const (
	TargetField  = "field"
	TargetCommon = "common"
)

// Error is the abstract error model
type Error struct {
	Title   string      `json:"title"`
	Details string      `json:"details"`
	Code    string      `json:"code"`
	Source  string      `json:"source,omitempty"`
	Target  string      `json:"target"`
	Meta    interface{} `json:"meta,omitempty"`
}

// NewError creates a new Error instance.
func NewError() *Error {
	return new(Error)
}

func NewCommonErrorByCode(code, details string) *Error {
	return NewCommonError().SetCode(code).SetTitleByCode(code).SetDetails(details)
}

// NewCommonError creates a new Error type of common instance.
func NewCommonError() *Error {
	return &Error{
		Target: TargetCommon,
	}
}

// NewFieldError creates a new Error type of field instance.
func NewFieldError() *Error {
	return &Error{
		Target: TargetField,
	}
}

// SetTitle sets the given title for the error object.
func (e *Error) SetTitle(title string) *Error {
	e.Title = title
	return e
}

// SetTitleByStatus sets the title for the error object by status code.
func (e *Error) SetTitleByStatus(status int) *Error {
	e.Title = http.StatusText(status)
	return e
}

// SetTitleByCode sets the title for the error object by code.
func (e *Error) SetTitleByCode(code string) *Error {
	if status, ok := statusCodes[code]; ok {
		e.SetTitleByStatus(status)
	}
	return e
}

// SetDetails sets the given details for the error object.
func (e *Error) SetDetails(details string) *Error {
	e.Details = details
	return e
}

// SetCode sets the code for the error object.
func (e *Error) SetCode(code string) *Error {
	e.Code = code
	return e
}

// SetSourсe sets the sourсe for the error object.
func (e *Error) SetSourсe(source string) *Error {
	e.Source = source
	return e
}

// SetTarget sets the sourсe for the error object.
func (e *Error) SetTarget(target string) *Error {
	e.Target = target
	return e
}

// ApplyCode sets the code and title by code for the error object.
func (e *Error) ApplyCode(code string) *Error {
	e.SetCode(code).SetTitleByCode(code)
	return e
}

func GetTypedError(code string) *errors.PublicError {
	return &errors.PublicError{
		Code:       code,
		HttpStatus: statusCodes[code],
	}
}
