package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"services.shared/apperror"
	"services.shared/logger"
)

type Responder interface {
	Success(c *gin.Context, data interface{})
	Error(c *gin.Context, err error)
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
func (r *responder) Success(c *gin.Context, data interface{}) {
	payload := successResponse{data}
	r.response(c, payload)
}

type successResponse struct {
	Data interface{} `json:"data"`
}

// Error responses error to client
func (r *responder) Error(c *gin.Context, err error) {
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

func (r *responder) response(c *gin.Context, payload interface{}) {
	c.JSON(http.StatusOK, payload)
}
