package errors

// Error codes

const (
	// The request was unacceptable, often due to missing a required parameter.
	CodeInvalidRequestError = "invalid_request"
	CodeInternalError       = "internal_error"
)

var codeText = map[string]string{
	CodeInvalidRequestError: "invalid request",
	CodeInternalError:       "internal error",
}

// CodeText returns a text for the error code. It returns the empty
// string if the code is unknown.
func CodeText(code string) string {
	return codeText[code]
}
