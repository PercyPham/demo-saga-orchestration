package msg_test

import (
	"services.shared/saga/msg"
	"testing"
)

func TestNewEvent(t *testing.T) {
	meta := msg.EventMeta{
		Topic:         "topic",
		CorrelationID: "correlationID",
		Type:          "eventType",
	}
	event := msg.NewEvent(meta, "payload")

	tests := []struct {
		expected string
		got      string
	}{
		{"topic", event.Topic()},
		{"correlationID", event.CorrelationID()},
		{"eventType", event.Type()},
		{"payload", event.Payload()},
	}

	for _, test := range tests {
		if test.got != test.expected {
			t.Errorf("expected %s, got %s", test.expected, test.got)
		}
	}
}

func TestValidateEvent(t *testing.T) {
	headers := map[string]string{
		msg.HeaderMessageType: msg.TypeEvent,
		msg.HeaderEventType:   "eventType",
	}
	message := msg.NewMessage(headers, "")
	event, err := msg.ValidateEvent(message)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if event == nil {
		t.Error("expected event, received nil")
	} else if event.Type() != "eventType" {
		t.Errorf("expected eventType, got %s", event.Type())
	}
}
