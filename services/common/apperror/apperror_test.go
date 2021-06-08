package apperror_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"services.common/apperror"
)

func TestNewAppError(t *testing.T) {
	appErr := apperror.New(apperror.NotFound, "message")

	assert.Equal(t, apperror.NotFound, appErr.Code())
	assert.Equal(t, "message", appErr.Message())
}

func TestWrapError(t *testing.T) {
	err := errors.New("origin error message")
	wrappedErr := apperror.Wrap(err, apperror.BadRequest, "wrapped error message", "wrapped log")

	assert.Equal(t, apperror.BadRequest, wrappedErr.Code())
	assert.Equal(t, "wrapped error message", wrappedErr.Message())
}

func TestWrapAppError(t *testing.T) {
	appErr := apperror.New(apperror.NotFound, "not found message")
	wrappedAppErr := apperror.Wrap(appErr, apperror.BadRequest, "bad request message", "bad request log")

	assert.Equal(t, apperror.BadRequest, wrappedAppErr.Code())
	assert.Equal(t, "bad request message", wrappedAppErr.Message())
}

func TestRootErrorOfNil(t *testing.T) {
	rootErr := apperror.RootError(nil)
	assert.Equal(t, nil, rootErr)
}

func TestRootErrorOfErr(t *testing.T) {
	err := errors.New("message")
	rootErr := apperror.RootError(err)

	assert.Equal(t, err, rootErr)

	wrappedErr1 := apperror.Wrap(err, apperror.BadRequest, "wrapped message 1", "wrapped log 1")
	wrappedErr2 := apperror.Wrap(wrappedErr1, apperror.BadRequest, "wrapped message 2", "wrapped log 2")
	rootErr = apperror.RootError(wrappedErr2)

	assert.Equal(t, err, rootErr)
}

func TestRootErrorOfAppErr(t *testing.T) {
	appErr := apperror.New(apperror.BadRequest, "app err message")
	rootErr := apperror.RootError(appErr)

	assert.Equal(t, appErr, rootErr)

	wrappedErr1 := apperror.Wrap(appErr, apperror.BadRequest, "wrapped message 1", "wrapped log 1")
	wrappedErr2 := apperror.Wrap(wrappedErr1, apperror.BadRequest, "wrapped message 2", "wrapped log 2")
	rootErr = apperror.RootError(wrappedErr2)

	assert.Equal(t, appErr, rootErr)
}

func TestAppErrorErrorMessage(t *testing.T) {
	appErr := apperror.New(apperror.NotFound, "message")

	assert.Equal(t, "message", appErr.Error())
}

func TestWrappedAppErrorErrorMessage(t *testing.T) {
	appErr := apperror.New(apperror.NotFound, "message")
	wrappedAppErr1 := apperror.Wrap(appErr, apperror.BadRequest, "wrapped message 1", "wrapped log 1")
	wrappedAppErr2 := apperror.Wrap(wrappedAppErr1, apperror.BadRequest, "wrapped message 2", "wrapped log 2")

	assert.Equal(t, "message", wrappedAppErr2.Error())
	assert.Equal(t, "wrapped message 2\nCaused by: wrapped log 2\nCaused by: wrapped log 1\nCaused by: message", wrappedAppErr2.Trace())
}

func TestWrappedErrorErrorMessage(t *testing.T) {
	err := errors.New("message")
	wrappedAppErr1 := apperror.Wrap(err, apperror.BadRequest, "wrapped message 1", "wrapped log 1")
	wrappedAppErr2 := apperror.Wrap(wrappedAppErr1, apperror.BadRequest, "wrapped message 2", "wrapped log 2")

	assert.Equal(t, "message", wrappedAppErr2.Error())
	assert.Equal(t, "wrapped message 2\nCaused by: wrapped log 2\nCaused by: wrapped log 1\nCaused by: message", wrappedAppErr2.Trace())
}

func TestErrWithLog(t *testing.T) {
	err := errors.New("message")
	appErr := apperror.WithLog(err, "log")

	assert.Equal(t, apperror.InternalServerError, appErr.Code())
	assert.Equal(t, "internal server error", appErr.Message())

	assert.Equal(t, err, apperror.RootError(appErr))

	assert.Equal(t, "message", appErr.Error())
	assert.Equal(t, "internal server error\nCaused by: log\nCaused by: message", appErr.Trace())
}

func TestAppErrWithLog(t *testing.T) {
	appErr := apperror.New(apperror.NotFound, "message")
	withLogAppErr := apperror.WithLog(appErr, "log")

	assert.Equal(t, apperror.NotFound, withLogAppErr.Code())
	assert.Equal(t, "message", withLogAppErr.Message())

	assert.Equal(t, appErr, apperror.RootError(withLogAppErr))

	assert.Equal(t, "message", withLogAppErr.Error())
	assert.Equal(t, "message\nCaused by: log\nCaused by: message", withLogAppErr.Trace())
}
