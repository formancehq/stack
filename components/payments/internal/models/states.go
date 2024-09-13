package models

import (
	"encoding/json"
)

type State struct {
	ID          StateID
	ConnectorID ConnectorID
	State       json.RawMessage
}

func (s State) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID          string          `json:"id"`
		ConnectorID string          `json:"connectorID"`
		State       json.RawMessage `json:"state"`
	}{
		ID:          s.ID.String(),
		ConnectorID: s.ConnectorID.String(),
		State:       s.State,
	})
}

func (s *State) UnmarshalJSON(data []byte) error {
	var aux struct {
		ID          string          `json:"id"`
		ConnectorID string          `json:"connectorID"`
		State       json.RawMessage `json:"state"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	id, err := StateIDFromString(aux.ID)
	if err != nil {
		return err
	}

	connectorID, err := ConnectorIDFromString(aux.ConnectorID)
	if err != nil {
		return err
	}

	s.ID = *id
	s.ConnectorID = connectorID
	s.State = aux.State

	return nil
}
