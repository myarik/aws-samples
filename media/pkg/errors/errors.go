package errors

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

// The wrapper simplify processing HTTP errors
// Example:
//			lerror.WrapWithStatusCode(err, http.StatusInternalServerError, "cannot decode a request")
//	 		lerror.WrapCode(err, lerror.CodeInvalidRequestError, "invalid request")

type restError struct {
	Type string `json:"type"`
	Msg  string `json:"message"`
}

// Error returns a error msg
func (e *restError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Msg)
}

// WrapStatusCode returns APIGatewayProxyResponse with a status code and error message
func WrapStatusCode(httpStatus int, message string) (events.APIGatewayProxyResponse, error) {
	if message == "" {
		message = http.StatusText(httpStatus)
	}

	body, err := json.Marshal(&restError{
		Type: CodeInternalError,
		Msg:  message,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      http.StatusInternalServerError,
			Body:            "cannot marshal error",
			IsBase64Encoded: false,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      httpStatus,
		Body:            string(body),
		IsBase64Encoded: false,
	}, nil
}

// WrapCode returns APIGatewayProxyResponse with the 400 status code and error message
func WrapCode(code, message string) (events.APIGatewayProxyResponse, error) {
	if message == "" {
		message = CodeText(code)
	}

	apiErr := &restError{
		Type: code,
		Msg:  message,
	}
	body, err := json.Marshal(apiErr)
	if err != nil {
		return WrapStatusCode(http.StatusInternalServerError, "cannot marshal error")
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusBadRequest,
		Body:            string(body),
		IsBase64Encoded: false,
	}, nil
}
