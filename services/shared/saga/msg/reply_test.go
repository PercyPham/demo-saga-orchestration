package msg_test

import (
	"services.shared/saga/msg"
	"testing"
)

func TestNewReply(t *testing.T) {
	meta := msg.ReplyMeta{
		Command: "cmdID",
		Success: true,
		Type:    "eventType",
	}
	reply := msg.NewReply(meta, "payload")

	tests := []struct {
		expected interface{}
		got      interface{}
	}{
		{"cmdID", reply.Command()},
		{true, reply.Success()},
		{"eventType", reply.Type()},
		{"payload", reply.Payload()},
	}

	for _, test := range tests {
		if test.got != test.expected {
			t.Errorf("expected %s, got %s", test.expected, test.got)
		}
	}
}

func TestValidateReply(t *testing.T) {
	headers := map[string]string{
		msg.HeaderMessageType: msg.TypeReply,
		msg.HeaderReplyType:   "replyType",
	}
	message := msg.NewMessage(headers, "")
	reply, err := msg.ValidateReply(message)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if reply == nil {
		t.Error("expected reply, received nil")
	} else if reply.Type() != "replyType" {
		t.Errorf("expected replyType, got %s", reply.Type())
	}
}
