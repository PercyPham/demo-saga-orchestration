package rest_response

import (
	"net/http"
	"services.shared/apperror"
	"services.shared/logger"
)

// JSONResponder is interface for REST api implementation to response JSON
//  use *gin.Context
type JSONResponder interface {
	// JSON serializes the given struct as JSON into the response body. It also sets the Content-Type as "application/json".
	JSON(code int, obj interface{})
}

type Responder interface {
	Success(c JSONResponder, data interface{})
	Error(c JSONResponder, err error)
	SetLogTrace(logTrace bool)
}

func New(log logger.Logger) Responder {
	return &responder{log, false}
}

type responder struct {
	log      logger.Logger
	logTrace bool
}

func (r *responder) SetLogTrace(v bool) {
	r.logTrace = v
}

// Success responses success data to client
func (r *responder) Success(c JSONResponder, data interface{}) {
	payload := successResponse{data}
	r.response(c, payload)
}

type successResponse struct {
	Data interface{} `json:"data"`
}

// Error responses error to client
func (r *responder) Error(c JSONResponder, err error) {
	var payload errorResponse

	appErr, ok := err.(*apperror.AppError)
	if ok {
		if r.logTrace || appErr.Code() == apperror.InternalServerError {
			r.log.Error(appErr.Trace())
		}
		payload = newErrorResponse(int(appErr.Code()), apperror.StatusText(appErr.Code()), appErr.Message())
		r.response(c, payload)
		return
	}

	r.log.Error(apperror.WithLog(err, "").Trace())
	payload = newErrorResponse(
		int(apperror.InternalServerError),
		apperror.StatusText(apperror.InternalServerError),
		"internal server error",
	)

	r.response(c, payload)
}

func newErrorResponse(code int, status string, message string) errorResponse {
	errContent := errorContent{code, status, message}
	return errorResponse{errContent}
}

type errorResponse struct {
	Error errorContent `json:"error"`
}

type errorContent struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (r *responder) response(c JSONResponder, payload interface{}) {
	c.JSON(http.StatusOK, payload)
}
