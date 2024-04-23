package messaging

import (
	"testing"
)

func TestMessage_ToMap(t *testing.T) {
	// Assuming the Data field is a map[string]interface{}
	// interface {}(map[string]interface {}) ["message_id": "9e2c0fb2-64b6-4176-9d90-184fab31e95c", "event_id": "dfdfecce-9d5a-41d4-be8d-2291a26a2578", "error": nil, "status": "Processed", ]
	m := Message{
		OperationID: "9e2c0fb2-64b6-4176-9d90-184fab31e95c",
		Error:       nil,
		Data: map[string]interface{}{
			"message_id": "9e2c0fb2-64b6-4176-9d90-184fab31e95c",
			"event_id":   "dfdfecce-9d5a-41d4-be8d-2291a26a2578",
			"error":      nil,
			"status":     "Processed",
		},
	}

	// Try to convert the Data field to a map[string]interface{}
	dataMap, err := m.ToMap()
	if err != nil {
		t.Errorf("Message.ToMap() error = %v", err)
		return
	}

	// Check if the Data field was converted to a map[string]interface{}
	if len(dataMap) == 0 {
		t.Errorf("Message.ToMap() = %v, want map[string]interface{}", dataMap)
	}
}
