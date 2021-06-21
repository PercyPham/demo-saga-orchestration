package msg_test

import (
	"services.shared/saga/msg"
	"testing"
)

func TestNewMessage(t *testing.T) {
	headers := map[string]string{
		msg.HeaderMessageID:   "id",
		msg.HeaderMessageType: "type",
		"key":                 "val",
	}
	message := msg.NewMessage(headers, "payload")

	tests := []struct {
		expected string
		got      string
	}{
		{"id", message.ID()},
		{"type", message.Type()},
		{"val", message.GetHeader("key")},
		{"payload", message.Payload()},
	}

	for _, test := range tests {
		if test.got != test.expected {
			t.Errorf("expected %s, got %s", test.expected, test.got)
		}
	}
}

func TestSetInfo(t *testing.T) {
	message := msg.NewMessage(nil, "")
	message.SetID("id")
	message.SetHeader("key", "val")
	message.SetPayload("payload")

	tests := []struct {
		expected string
		got      string
	}{
		{"id", message.ID()},
		{"val", message.GetHeader("key")},
		{"payload", message.Payload()},
	}

	for _, test := range tests {
		if test.got != test.expected {
			t.Errorf("expected %s, got %s", test.expected, test.got)
		}
	}
}

func TestMarshalMessage(t *testing.T) {
	headers := map[string]string{
		"key": "val",
	}
	message := msg.NewMessage(headers, "payload")
	jsonMsg, err := msg.Marshal(message)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	expectedJsonMsg := "{\"headers\":{\"key\":\"val\"},\"payload\":\"payload\"}"
	if string(jsonMsg) != expectedJsonMsg {
		t.Errorf("expected %s, got %s", expectedJsonMsg, string(jsonMsg))
	}
}

func TestUnmarshalMessage(t *testing.T) {
	raw := []byte("{\"headers\":{\"key\":\"val\"},\"payload\":\"payload\"}")
	message, err := msg.Unmarshal(raw)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if message == nil {
		t.Errorf("expected message, got nil")
	} else {
		if message.GetHeader("key") != "val" {
			t.Errorf("expected val, got %v", message.GetHeader("key"))
		}
		if message.Payload() != "payload" {
			t.Errorf("expected payload, got %s", message.Payload())
		}
	}
}
