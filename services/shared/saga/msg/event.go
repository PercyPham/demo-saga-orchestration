package msg

import "errors"

const (
	HeaderEventTopic         = "event_topic"
	HeaderEventCorrelationID = "event_correlation_id"
	HeaderEventType          = "event_type"
)

type Event interface {
	Message

	Topic() string
	CorrelationID() string
	Type() string

	SetCorrelationID(id string)
}

type EventMeta struct {
	Topic         string
	CorrelationID string
	Type          string
}

func NewEvent(meta EventMeta, payload string) Event {
	headers := map[string]string{
		HeaderMessageType:        TypeEvent,
		HeaderEventTopic:         meta.Topic,
		HeaderEventCorrelationID: meta.CorrelationID,
		HeaderEventType:          meta.Type,
	}
	message := NewMessage(headers, payload)
	return &event{message}
}

type event struct {
	Message
}

func (e *event) Topic() string              { return e.Header(HeaderEventTopic) }
func (e *event) CorrelationID() string      { return e.Header(HeaderEventCorrelationID) }
func (e *event) Type() string               { return e.Header(HeaderEventType) }
func (e *event) SetCorrelationID(id string) { e.SetHeader(HeaderEventCorrelationID, id) }

func ValidateEvent(message Message) (Event, error) {
	if message.Type() != TypeEvent {
		return nil, errors.New("invalid event message type")
	}
	return &event{message}, nil
}
