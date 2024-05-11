package messaging

import (
	"encoding/json"
	"errors"
)

type Message interface {
	GetError() error
	GetStatus() string
	GetCommand() string
	GetData() []byte
	Deserialize(message []byte) error
	Serialize() ([]byte, error)
}

// MessageImpl implements the Message interface
type MessageImpl struct {
	Error   error  `json:"error"`
	Status  string `json:"status"`
	Command string `json:"command"`
	Data    []byte `json:"data"`
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
func (m *MessageImpl) GetData() []byte {
	return m.Data
}

// Deserializes a byte slice into Message
func (m *MessageImpl) Deserialize(message []byte) error {
	if m == nil {
		err := errors.New("deserialize::message is nil")
		return err
	}

	err := json.Unmarshal(message, m)
	return err
}

// Serializes a Message into a byte slice
func (m *MessageImpl) Serialize() ([]byte, error) {
	if m == nil {
		err := errors.New("serialize::message is nil")
		return nil, err
	}

	data, err := json.Marshal(m)
	return data, err
}

// NewMessage creates a new message
func NewMessage(error error, status string, command string, data []byte) Message {
	msg := &MessageImpl{
		Error:   error,
		Status:  status,
		Command: command,
		Data:    data,
	}

	return msg
}
