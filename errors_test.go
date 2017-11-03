package chargehound_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/chargehound/chargehound-go"
)

var errorTests = []struct {
	error     string
	code      int
	errorType chargehound.ErrorType
	message   string
}{
	{
		"{\"error\": { \"status\": 404, \"message\": \"A dispute with id 'puppy' was not found\"}}",
		404,
		chargehound.NotFoundError,
		"Not Found: A dispute with id 'puppy' was not found",
	},
	{
		"{\"error\": { \"status\": 400, \"message\": \"Wrong param\"}}",
		400,
		chargehound.BadRequestError,
		"Bad Request: Wrong param",
	},
	{
		"{\"error\": { \"status\": 401, \"message\": \"No user\"}}",
		401,
		chargehound.UnauthorizedError,
		"Unauthorized: No user",
	},
	{
		"{\"error\": { \"status\": 403, \"message\": \"Wrong user\"}}",
		403,
		chargehound.ForbiddenError,
		"Forbidden: Wrong user",
	},
	{
		"{\"error\": { \"status\": 500, \"message\": \"Server error\"}}",
		500,
		chargehound.InternalServerError,
		"Server Error: Server error",
	},
}

func TestErrors(t *testing.T) {
	ch := chargehound.New("api_key", nil)

	for _, test := range errorTests {
		error := test.error
		code := test.code
		errorType := test.errorType
		message := test.message

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, error, code)
		}))
		defer ts.Close()

		url, err := url.Parse(ts.URL)
		if err != nil {
			t.Error(err)
		}

		ch.Host = url.Host
		ch.Protocol = url.Scheme + "://"

		_, err = ch.Disputes.Retrieve(&chargehound.RetrieveDisputeParams{ID: "puppy"})
		if err == nil {
			t.Error(err)
		}

		if message != err.Error() {
			t.Error("Expected error: ", message)
		}

		chErr := err.(chargehound.Error)

		if code != chErr.StatusCode() {
			t.Error("Expected status: ", code)
		}

		if errorType != chErr.Type() {
			t.Error("Expected type: ", errorType)
		}
	}
}
