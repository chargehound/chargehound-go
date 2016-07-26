package chargehound

import (
	"encoding/json"
	"net/http"
)

type ErrorType string

const (
	BadRequestError     = ErrorType("Bad Request")
	UnauthorizedError   = ErrorType("Unauthorized")
	ForbiddenError      = ErrorType("Forbidden")
	NotFoundError       = ErrorType("Not Found")
	InternalServerError = ErrorType("Server Error")
	GenericError        = ErrorType("Error")
)

type Error interface {
	Error() string
	StatusCode() int
	Type() ErrorType
}

type errorResponse struct {
	Status    int
	ErrorType ErrorType
	Message   string
}

type errorJSON struct {
	Url      string
	Livemode bool
	Error    errorResponse
}

func responseToError(res *http.Response) error {
	var errRes errorJSON
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(&errRes)
	if err != nil {
		return err
	}

	switch errRes.Error.Status {
	case 400:
		errRes.Error.ErrorType = BadRequestError
	case 401:
		errRes.Error.ErrorType = UnauthorizedError
	case 403:
		errRes.Error.ErrorType = ForbiddenError
	case 404:
		errRes.Error.ErrorType = NotFoundError
	case 500:
		errRes.Error.ErrorType = InternalServerError
	default:
		errRes.Error.ErrorType = GenericError
	}

	return &errRes.Error
}

func (e *errorResponse) Error() string {
	return string(e.ErrorType) + ": " + e.Message
}

func (e *errorResponse) StatusCode() int {
	return e.Status
}

func (e *errorResponse) Type() ErrorType {
	return e.ErrorType
}
