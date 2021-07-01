package apperr

import "fmt"

type AppErr interface {
	Error() string
	Code() int
	WithCode(code int) AppErr
	PublicMessage() string
	WithPublicMessage(string) AppErr
	WithPublicMessagef(format string, args ...interface{}) AppErr
}

func New(message string) AppErr {
	return &appErr{msg: message, code: InternalServerError}
}

func Wrap(err error, message string) AppErr {
	return wrap(err, message)
}

func Wrapf(err error, format string, args ...interface{}) AppErr {
	message := fmt.Sprintf(format, args...)
	return wrap(err, message)
}

func wrap(err error, message string) AppErr {
	if err == nil {
		return nil
	}
	code := InternalServerError
	if aErr, ok := err.(AppErr); ok {
		code = aErr.Code()
	}
	return &appErr{cause: err, msg: message, code: code}
}

type appErr struct {
	cause         error
	msg           string
	code          int
	publicMsg string
}

func (e *appErr) Error() string {
	if e.cause == nil {
		return e.msg
	}
	return e.msg + ": " + e.cause.Error()
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *appErr) Unwrap() error { return e.cause }

func (e *appErr) Code() int { return e.code }

func (e *appErr) WithCode(code int) AppErr {
	e.code = code
	return e
}

func (e *appErr) PublicMessage() string {
	if e.publicMsg != "" {
		return e.publicMsg
	}
	if e.cause == nil {
		return "internal server error"
	}
	if aErr, ok := e.cause.(AppErr); ok {
		return aErr.PublicMessage()
	}
	return "internal server error"
}

func (e *appErr) WithPublicMessage(m string) AppErr {
	e.publicMsg = m
	return e
}

func (e *appErr) WithPublicMessagef(format string, args ...interface{}) AppErr {
	e.publicMsg = fmt.Sprintf(format, args...)
	return e
}
