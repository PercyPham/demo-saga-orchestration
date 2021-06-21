package saga

type Saga struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	CurrentStep   int    `json:"current_step"`
	LastCommandID string `json:"last_command_id"`
	EndState      bool   `json:"end_state"`
	Compensating  bool   `json:"compensating"`
	Data          []byte `json:"data"`
}

// NewSaga return new saga instance
func NewSaga(sagaType string, data []byte) *Saga {
	return &Saga{
		Type: sagaType,
		Data: data,
	}
}
