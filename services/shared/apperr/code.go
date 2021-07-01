package apperr

// AppError's code
// For error codes UNDER 1000, these error codes are borrow-from/based-on HTTP status code.
//  See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
// For error codes from 1000 and above, these are app specific error codes
const (
	BadRequest          = 400
	Unauthorized        = 401
	Forbidden           = 403
	NotFound            = 404
	NotAcceptable       = 406
	UnprocessableEntity = 422

	InternalServerError = 500

	InvalidEvent   = 1000
	InvalidCommand = 1001
	InvalidReply   = 1002
)

var statusText = map[int]string{
	BadRequest:          "Bad Request",
	Unauthorized:        "Unauthorized",
	Forbidden:           "Forbidden",
	NotFound:            "Not Found",
	NotAcceptable:       "Not Acceptable",
	UnprocessableEntity: "Unprocessable Entity",

	InternalServerError: "Internal Server Error",

	InvalidEvent:   "Invalid Event",
	InvalidCommand: "Invalid Command",
	InvalidReply:   "Invalid Reply",
}

// StatusText returns a text for the app error code. It returns the empty
// string if the code is unknown.
func StatusText(code int) string {
	return statusText[code]
}
