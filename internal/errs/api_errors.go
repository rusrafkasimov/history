package errs

import (
	"net/http"
)

type ApiErr struct {
	message string
	status  int
	err     string
}

func (a *ApiErr) String() string {
	return a.message
}

func (a *ApiErr) Error() string {
	return a.err
}

func (a *ApiErr) ErrorMessage() string {
	return a.message
}

func (a *ApiErr) ErrorCode() int {
	return a.status
}


func NewBadRequestError(message string) *ApiErr {
	return &ApiErr{
		message: message,
		status:  http.StatusBadRequest,
		err:     "bad_request",
	}
}

func NewNotFoundError(message string) *ApiErr {
	return &ApiErr{
		message: message,
		status:  http.StatusNotFound,
		err:     "not_found_error",
	}
}

func NewInternalServerError(message string) *ApiErr {
	return &ApiErr{
		message: message,
		status:  http.StatusInternalServerError,
		err:     "internal_server_error",
	}
}

func NewUnauthorizedError(message string) *ApiErr {
	return &ApiErr{
		message: message,
		status:  http.StatusUnauthorized,
		err:     "unauthorized_error",
	}
}

func NewUnprocessableEntityError(message string) *ApiErr {
	return &ApiErr{
		message: message,
		status:  http.StatusUnprocessableEntity,
		err:     "unprocessable_entity_error",
	}
}

func NewNoContentError(message string) *ApiErr {
	return &ApiErr{
		message: message,
		status:  http.StatusNoContent,
		err:     "not_content_status",
	}
}

func NewTeapotError(message string) *ApiErr {
	return &ApiErr{
		message: message,
		status:  http.StatusTeapot,
		err:     "teapot_error",
	}
}