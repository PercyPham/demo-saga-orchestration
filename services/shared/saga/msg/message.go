package msg

import (
	"encoding/json"
	"errors"
)

const (
	HeaderMessageID   = "id"
	HeaderMessageType = "type"

	TypeCommand = "COMMAND"
	TypeReply   = "REPLY"
	TypeEvent   = "EVENT"
)

type Message interface {
	ID() string
	Type() string

	Headers() map[string]string
	Payload() string

	Header(key string) string

	SetID(id string)
	SetHeaders(map[string]string)
	SetHeader(key, val string)
	SetPayload(payload string)
}

type message struct {
	headers map[string]string
	payload string
}

func NewMessage(headers map[string]string, payload string) Message {
	if headers == nil {
		headers = make(map[string]string)
	}
	return &message{headers, payload}
}

func (m *message) ID() string                           { return m.headers[HeaderMessageID] }
func (m *message) Type() string                         { return m.headers[HeaderMessageType] }
func (m *message) Headers() map[string]string           { return m.headers }
func (m *message) Payload() string                      { return m.payload }
func (m *message) Header(key string) string             { return m.headers[key] }
func (m *message) SetID(id string)                      { m.headers[HeaderMessageID] = id }
func (m *message) SetHeaders(headers map[string]string) { m.headers = headers }
func (m *message) SetHeader(key, val string)            { m.headers[key] = val }
func (m *message) SetPayload(payload string)            { m.payload = payload }

// Marshal returns the JSON encoding of message
func Marshal(m Message) ([]byte, error) {
	jMsg := jsonMessage{
		Headers: m.Headers(),
		Payload: m.Payload(),
	}
	return json.Marshal(jMsg)
}

// Unmarshal parses the JSON-encoded message and returns message
func Unmarshal(jsonEncodedMsg []byte) (Message, error) {
	jMsg := jsonMessage{}
	err := json.Unmarshal(jsonEncodedMsg, &jMsg)
	if err != nil {
		return nil, errors.New("cannot unmarshal message from jsonEncodedMsg, got error: " + err.Error())
	}
	message := NewMessage(jMsg.Headers, jMsg.Payload)
	return message, nil
}

type jsonMessage struct {
	Headers map[string]string `json:"headers,omitempty"`
	Payload string            `json:"payload,omitempty"`
}
