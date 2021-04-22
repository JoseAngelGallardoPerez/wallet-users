package responses

// Response is the abstract response model
type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Errors []*Error    `json:"errors,omitempty"`
}

// NewResponse creates a new Response instance.
func NewResponse() *Response {
	return new(Response)
}

// SetStatus sets the status for the response object.
func (r *Response) SetStatus(status int) *Response {
	r.Status = status
	return r
}

// SetStatusByCode sets the status for the response object by code.
func (r *Response) SetStatusByCode(code string) *Response {
	if status, ok := statusCodes[code]; ok {
		r.Status = status
	}
	return r
}

// SetData sets the data for the response object.
func (r *Response) SetData(data interface{}) *Response {
	r.Data = data
	return r
}

// SetErrors sets list of errors for the response object.
func (r *Response) SetErrors(errs []*Error) *Response {
	r.Errors = errs
	return r
}

// AddError adds error into the response object.
func (r *Response) AddError(err *Error) *Response {
	r.Errors = append(r.Errors, err)
	return r
}
