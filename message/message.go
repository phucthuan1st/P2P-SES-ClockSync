package message

import (
	"encoding/json"
	"p2p-ses-clocksync/vectorclock"
)

type Payload struct {
	Name  string
	Clock []vectorclock.ClockEntry
}

type Message struct {
	Source    string
	Dest      string
	Content   string
	Timestamp []vectorclock.ClockEntry
	Payloads  []Payload
}

// Serialize serializes a Message struct to a JSON string
func (m *Message) Serialize() (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Deserialize deserializes a JSON string into a Message struct
func (m *Message) Deserialize(serialized string) error {
	err := json.Unmarshal([]byte(serialized), m)
	if err != nil {
		return err
	}
	return nil
}
