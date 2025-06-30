package models

import (
	"encoding/json"
	"time"
)

type UserActivityEvent struct {
	UserID    string                 `json:"userId"`
	EventType string                 `json:"eventType"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

func (e *UserActivityEvent) UnmarshalJSON(data []byte) error {
	type Alias UserActivityEvent
	aux := &struct {
		Timestamp string `json:"timestamp"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	parsedTime, err := time.Parse(time.RFC3339Nano, aux.Timestamp)
	if err != nil {
		parsedTime, err = time.Parse(time.RFC3339, aux.Timestamp)
		if err != nil {
			return err
		}
	}
	e.Timestamp = parsedTime
	return nil
}

func (e UserActivityEvent) DataToJson() ([]byte, error) {
	if e.Data == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(e.Data)
}
