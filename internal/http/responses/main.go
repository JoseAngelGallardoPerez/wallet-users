package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"

	errorsPkg "github.com/Confialink/wallet-pkg-errors"
)

// ResponseHandler is interface for response functionality
type ResponseHandler interface {
	SuccessResponse(ctx *gin.Context, status int, data interface{})
	ErrorResponse(ctx *gin.Context, status int, details string)
	Error(ctx *gin.Context, code, details string)
	SetError(ctx *gin.Context, err *Error)
	Errors(ctx *gin.Context, status int, errs []*Error)
	CustomErrorResponse(ctx *gin.Context, status int, title, details string)
	OkResponse(ctx *gin.Context, data interface{})
	ForbiddenResponse(ctx *gin.Context)
	Forbidden(ctx *gin.Context)
	NotFoundResponse(ctx *gin.Context)
	NotFound(ctx *gin.Context)
	ValidatorErrorResponse(ctx *gin.Context, code string, err error)
}

// ResponseService is an empty structure
type ResponseService struct{}

// NewResponseService creates a new ResponseService instance.
func NewResponseService() ResponseHandler {
	return &ResponseService{}
}

// SuccessResponse returns a success response
func (res *ResponseService) SuccessResponse(ctx *gin.Context, status int, data interface{}) {
	r := NewResponse().SetStatus(status).SetData(data)
	ctx.JSON(status, r)
}

// ErrorResponse returns a error response
func (res *ResponseService) ErrorResponse(ctx *gin.Context, status int, details string) {
	e := NewError().SetTitleByStatus(status).SetDetails(details)
	r := NewResponse().SetStatus(status).AddError(e)
	ctx.AbortWithStatusJSON(status, r)
}

// Error returns a error response [ Error is new version of ErrorResponse ]
func (res *ResponseService) Error(ctx *gin.Context, code, details string) {
	e := NewError().SetCode(code).SetTitleByCode(code).SetTarget(TargetCommon).SetDetails(details)
	r := NewResponse().SetStatusByCode(code).AddError(e)
	ctx.AbortWithStatusJSON(r.Status, r)
}

func (res *ResponseService) SetError(ctx *gin.Context, err *Error) {
	r := NewResponse().SetStatusByCode(err.Code).AddError(err)
	ctx.AbortWithStatusJSON(r.Status, r)
}

// Errors returns a list of errors response [ Errors is new version of ErrorResponse ]
func (res *ResponseService) Errors(ctx *gin.Context, status int, errs []*Error) {
	r := NewResponse().SetStatus(status).SetErrors(errs)
	ctx.AbortWithStatusJSON(r.Status, r)
}

func (res *ResponseService) CustomErrorResponse(ctx *gin.Context, status int, title, details string) {
	e := NewError().SetTitle(title).SetDetails(details)
	r := NewResponse().SetStatus(status).AddError(e)
	ctx.AbortWithStatusJSON(status, r)
}

// OkResponse returns a "200 StatusOK" response
func (res *ResponseService) OkResponse(ctx *gin.Context, data interface{}) {
	status := http.StatusOK
	r := NewResponse().SetStatus(status).SetData(data)
	ctx.JSON(status, r)
}

// ForbiddenResponse returns a "403 StatusForbidden" response
func (res *ResponseService) ForbiddenResponse(ctx *gin.Context) {
	status := http.StatusForbidden
	e := NewError().SetTitleByStatus(status).SetDetails("User is not logged in.")
	r := NewResponse().SetStatus(status).AddError(e)
	ctx.AbortWithStatusJSON(status, r)
}

// Forbidden returns a "403 StatusForbidden" response [ Forbidden is newest version of ForbiddenResponse ]
func (res *ResponseService) Forbidden(ctx *gin.Context) {
	e := NewError().
		SetCode(Forbidden).
		SetTitleByCode(Forbidden).
		SetTarget(TargetCommon).
		SetDetails("User is not logged in.")
	r := NewResponse().SetStatusByCode(Forbidden).AddError(e)
	ctx.AbortWithStatusJSON(r.Status, r)
}

// NotFoundResponse returns a "404 StatusNotFound" response
func (res *ResponseService) NotFoundResponse(ctx *gin.Context) {
	status := http.StatusNotFound
	e := NewError().SetTitleByStatus(status).SetDetails("Not Found.")
	r := NewResponse().SetStatus(status).AddError(e)
	ctx.AbortWithStatusJSON(status, r)
}

// NotFound returns a "404 StatusNotFound" response [ NotFound is newest version of NotFoundResponse ]
func (res *ResponseService) NotFound(ctx *gin.Context) {
	e := NewError().SetCode(NotFound).SetTitleByCode(NotFound).SetTarget(TargetCommon).SetDetails("Not Found.")
	r := NewResponse().SetStatusByCode(NotFound).AddError(e)
	ctx.AbortWithStatusJSON(r.Status, r)
}

// ValidatorErrorResponse returns a error response
func (res *ResponseService) ValidatorErrorResponse(ctx *gin.Context, code string, err error) {
	if err, ok := err.(errorsPkg.TypedError); ok {
		_ = ctx.Error(err)
		return
	}

	errorsPkg.AddShouldBindError(ctx, err)
}
