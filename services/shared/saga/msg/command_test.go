package msg_test

import (
	"services.shared/saga/msg"
	"testing"
)

func TestNewCommand(t *testing.T) {
	meta := msg.CommandMeta{
		Type:        "commandType",
		Destination: "destination",
	}
	command := msg.NewCommand(meta, "payload")

	tests := []struct {
		expected string
		got      string
	}{
		{"commandType", command.Type()},
		{"destination", command.Destination()},
		{"payload", command.Payload()},
	}

	for _, test := range tests {
		if test.got != test.expected {
			t.Errorf("expected %s, got %s", test.expected, test.got)
		}
	}
}

func TestValidateCommand(t *testing.T) {
	headers := map[string]string{
		msg.HeaderMessageType: msg.TypeCommand,
		msg.HeaderCommandType: "commandType",
	}
	message := msg.NewMessage(headers, "")
	command, err := msg.ValidateCommand(message)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if command == nil {
		t.Error("expected command, received nil")
	} else if command.Type() != "commandType" {
		t.Errorf("expected commandType, got %s", command.Type())
	}
}
