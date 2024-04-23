package message

import "context"

// Interface for messaging systems
type MessagingSystem interface {
	Publish(ctx context.Context, data interface{}) error
	Subscribe(ctx context.Context) (<-chan Message, context.CancelFunc, error)
	Close(ctx context.Context) error
}

type Message interface {
	GetOperationID() string
	GetError() error
	GetStatus() string
	GetData() interface{}
}

// MessageImpl implements the Message interface
type MessageImpl struct {
	OperationID string      `json:"operation_id"`
	Error       error       `json:"error"`
	Status      string      `json:"status"`
	Data        interface{} `json:"data"`
}

// GetOperationID returns the operation ID
func (m *MessageImpl) GetOperationID() string {
	return m.OperationID
}

// GetError returns the error
func (m *MessageImpl) GetError() error {
	return m.Error
}

// GetStatus returns the status
func (m *MessageImpl) GetStatus() string {
	return m.Status
}

// GetData returns the data
func (m *MessageImpl) GetData() interface{} {
	return m.Data
}

// NewMessage creates a new message
func NewMessage(operationID string, error error, status string, data interface{}) Message {
	return &MessageImpl{
		OperationID: operationID,
		Error:       error,
		Status:      status,
		Data:        data,
	}
}
