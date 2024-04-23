package message

import (
	"context"
	"encoding/json"
	"errors"
)

// Interface for messaging systems
type MessagingSystem interface {
	Publish(ctx context.Context, data Message) error
	Subscribe(ctx context.Context) (<-chan Message, context.CancelFunc, error)
	Close(ctx context.Context) error
}

type Message interface {
	GetOperationID() string
	GetError() error
	GetStatus() string
	GetCommand() string
	GetData() interface{}
	SerializeData() ([]byte, error)
	DeserializeData(data []byte) error
}

// MessageImpl implements the Message interface
type MessageImpl struct {
	OperationID string      `json:"operation_id"`
	Error       error       `json:"error"`
	Status      string      `json:"status"`
	Command     string      `json:"command"`
	Data        interface{} `json:"-"`
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

// GetCommand returns the command
func (m *MessageImpl) GetCommand() string {
	return m.Command
}

// GetData returns the data
func (m *MessageImpl) GetData() interface{} {
	return m.Data
}

// SerializeData serializes the Data field to a byte slice
func (m *MessageImpl) SerializeData() ([]byte, error) {
	return json.Marshal(m.Data)
}

// DeserializeData deserializes a byte slice into the Data field
func (m *MessageImpl) DeserializeData(data []byte) error {
	if m.Data == nil {
		return errors.New("data field is nil")
	}
	return json.Unmarshal(data, &m.Data)
}

// NewMessage creates a new message with serialized data
func NewMessage(operationID string, error error, status string, command string, data interface{}) (Message, error) {
	msg := &MessageImpl{
		OperationID: operationID,
		Error:       error,
		Status:      status,
		Command:     command,
		Data:        data,
	}

	return msg, nil
}
