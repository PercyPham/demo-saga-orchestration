package apperror

type AppError struct {
	err     error
	code    int
	message string
	log     string
}

// New returns a new AppError
//
// AppError's code
// For error codes UNDER 1000, these error codes are borrow-from/based-on HTTP status code.
//  See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
// For error codes from 1000 and above, these are app specific error codes.
func New(code int, message string) *AppError {
	return &AppError{
		code:    code,
		message: message,
	}
}

// Wrap wraps input error with addition infos and return AppError
func Wrap(err error, code int, message, log string) *AppError {
	return &AppError{
		err:     err,
		code:    code,
		message: message,
		log:     log,
	}
}

// WithLog wraps error with addition log and return AppError
//  if input error is AppError, it will retain code, message and add log to returned AppError
//  if input error is plain error, it will wrap input with InternalServerError and add log to returned AppError
func WithLog(err error, log string) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return Wrap(err, appErr.code, appErr.message, log)
	}

	return Wrap(err, InternalServerError, "internal server error", log)
}

// RootError return most inner error of input error
//  if input error is nil, it will return nil
//  if input error is plain error, it will return input error
func RootError(err error) error {
	tempErr := err
	for {
		ae, ok := tempErr.(*AppError)
		if !ok {
			return tempErr
		}
		if ae.err == nil {
			return ae
		}
		tempErr = ae.err
	}
}

// Error returns root error message
func (e *AppError) Error() string {
	rootErr := RootError(e)
	if appErrRoot, ok := rootErr.(*AppError); ok {
		return appErrRoot.message
	}
	return rootErr.Error()
}

// Trace returns stack trace of error, including all logs and root error message
func (e *AppError) Trace() string {
	rootErr := RootError(e)
	logs := e.logs()
	if len(logs) == 0 {
		if appErrRoot, ok := rootErr.(*AppError); ok {
			return appErrRoot.message
		}
		return rootErr.Error()
	}
	s := e.Message()
	for _, log := range logs {
		s += "\nCaused by: " + log
	}
	s += "\nCaused by: " + rootErr.Error()
	return s
}

// logs return all logs from outer most to inner most error, excluding empty logs
func (e *AppError) logs() []string {
	logs := make([]string, 0)

	if e.log != "" {
		logs = []string{e.log}
	}

	err := e.err

	for err != nil {
		appErr, ok := err.(*AppError)
		if !ok {
			break
		}
		if appErr.log != "" {
			logs = append(logs, appErr.log)
		}
		err = appErr.err
	}

	return logs
}

func (e *AppError) Code() int {
	return e.code
}

func (e *AppError) Message() string {
	return e.message
}
