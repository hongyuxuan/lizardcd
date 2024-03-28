package errorx

import "net/http"

type LizardcdError struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type HttpErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewError(code int, message string, data interface{}) error {
	return &LizardcdError{Code: code, Message: message, Data: data}
}

func NewDefaultError(message string) error {
	return &LizardcdError{Code: http.StatusInternalServerError, Message: message}
}

func (e *LizardcdError) Error() string {
	return e.Message
}

func (e *LizardcdError) GetData() *HttpErrorResponse {
	return &HttpErrorResponse{
		Code:    e.Code,
		Message: e.Message,
		Data:    e.Data,
	}
}
