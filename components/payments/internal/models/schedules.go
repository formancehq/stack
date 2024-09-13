package models

import (
	"encoding/json"
	"time"
)

type Schedule struct {
	ID          string
	ConnectorID ConnectorID
	CreatedAt   time.Time
}

func (s Schedule) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID          string    `json:"id"`
		ConnectorID string    `json:"connectorID"`
		CreatedAt   time.Time `json:"createdAt"`
	}{
		ID:          s.ID,
		ConnectorID: s.ConnectorID.String(),
		CreatedAt:   s.CreatedAt,
	})
}

func (s *Schedule) UnmarshalJSON(data []byte) error {
	var aux struct {
		ID          string    `json:"id"`
		ConnectorID string    `json:"connectorID"`
		CreatedAt   time.Time `json:"createdAt"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	connectorID, err := ConnectorIDFromString(aux.ConnectorID)
	if err != nil {
		return err
	}

	s.ID = aux.ID
	s.ConnectorID = connectorID
	s.CreatedAt = aux.CreatedAt

	return nil
}
