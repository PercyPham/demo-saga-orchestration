package apperror

import "fmt"

type AppError interface {
	Error() string
	Code() int
	WithCode(code int) AppError
	PublicMessage() string
	WithPublicMessage(string) AppError
	WithPublicMessagef(format string, args ...interface{}) AppError
}

func New(message string) AppError {
	return &appErr{msg: message, code: InternalServerError}
}

func Newf(format string, args ...interface{}) AppError {
	message:= fmt.Sprintf(format, args...)
	return &appErr{msg: message, code: InternalServerError}
}

func Wrap(err error, message string) AppError {
	return wrap(err, message)
}

func Wrapf(err error, format string, args ...interface{}) AppError {
	message := fmt.Sprintf(format, args...)
	return wrap(err, message)
}

func wrap(err error, message string) AppError {
	if err == nil {
		return nil
	}
	code := InternalServerError
	if aErr, ok := err.(AppError); ok {
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

func (e *appErr) WithCode(code int) AppError {
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
	if aErr, ok := e.cause.(AppError); ok {
		return aErr.PublicMessage()
	}
	return "internal server error"
}

func (e *appErr) WithPublicMessage(m string) AppError {
	e.publicMsg = m
	return e
}

func (e *appErr) WithPublicMessagef(format string, args ...interface{}) AppError {
	e.publicMsg = fmt.Sprintf(format, args...)
	return e
}
