package apperror_test

import (
	"errors"
	"services.shared/apperror"
	"testing"
)

func TestNewAppErr(t *testing.T) {
	appErr := apperror.New("new error")
	if appErr.Error() != "new error" {
		t.Errorf("expected 'new error', got '%s'", appErr.Error())
	}
}

func TestNewAppErrf(t *testing.T) {
	appErr := apperror.Newf("new %s", "error")
	if appErr.Error() != "new error" {
		t.Errorf("expected 'new error', got '%s'", appErr.Error())
	}
}

func TestWrapNil(t *testing.T) {
	appErr := apperror.Wrap(nil, "error message")
	if appErr != nil {
		t.Errorf("expected nil, got '%v'", appErr)
	}
}

func TestWrapError(t *testing.T) {
	err := errors.New("error")
	appErr := apperror.Wrap(err, "wrap")
	if appErr.Error() != "wrap: error" {
		t.Errorf("expected 'wrap: error', got '%s'", appErr.Error())
	}
}

func TestWrapWithFormatMessage(t *testing.T) {
	err := errors.New("error")
	appErr := apperror.Wrapf(err, "formatted %s %v", "message", 0)
	expected := "formatted message 0: error"
	if appErr.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, appErr.Error())
	}
}

func TestWithoutCode(t *testing.T) {
	appErr := apperror.New("error")
	if appErr.Code() != apperror.InternalServerError {
		t.Errorf("expected %v, got %v", apperror.InternalServerError, appErr.Code())
	}
}

func TestWithCode(t *testing.T) {
	appErr := apperror.New("error").WithCode(1000)
	if appErr.Code() != 1000 {
		t.Errorf("expected %v, got %v", 1000, appErr.Code())
	}
}

func TestWithoutPublicMessage(t *testing.T) {
	appErr := apperror.New("error")
	if appErr.PublicMessage() != "internal server error" {
		t.Errorf("expected 'internal server error', got %s", appErr.PublicMessage())
	}
}

func TestWithPublicMessage(t *testing.T) {
	appErr := apperror.New("error").WithPublicMessage("public message")
	if appErr.PublicMessage() != "public message" {
		t.Errorf("expected 'public message', got %s", appErr.PublicMessage())
	}
}

func TestWithPublicMessagef(t *testing.T) {
	appErr := apperror.New("error").WithPublicMessagef("public %s", "message")
	expected := "public message"
	if appErr.PublicMessage() != expected {
		t.Errorf("expected '%s', got %s", expected, appErr.PublicMessage())
	}
}
