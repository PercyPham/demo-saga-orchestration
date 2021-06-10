package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"services.order/internal/common/config"
	"services.shared/apperror"
	"services.shared/logger"
)

type Responder interface {
	Success(c *gin.Context, data interface{})
	Error(c *gin.Context, err error)
}

func New(log logger.Logger) Responder {
	return &responder{log}
}

type responder struct {
	log logger.Logger
}

// Success responses success data to client
func (r responder) Success(c *gin.Context, data interface{}) {
	payload := successResponse{data}
	r.response(c, payload)
}

type successResponse struct {
	Data interface{} `json:"data"`
}

// Error responses error to client
func (r *responder) Error(c *gin.Context, err error) {
	var payload errorResponse
	isDev := config.App().ENV == "development"

	appErr, ok := err.(*apperror.AppError)
	if ok {
		if isDev || appErr.Code() == apperror.InternalServerError {
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
