package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Channel string

const (
	HookChannel    Channel = "HOOK_CHANNEL"
	AttemptChannel Channel = "ATTEMPT_CHANNEL"
)

type EventType string

const (
	NewWaitingAttemptType    EventType = "NEW_WAITING_ATTEMPT_TYPE"
	FlushWaitingAttemptType  EventType = "FLUSH_WAITING_ATTEMPT"
	FlushWaitingAttemptsType EventType = "FLUSH_WAITING_ATTEMPTS"
	AbortWaitingAttemptType  EventType = "ABORT_WAITING_ATTEMPT"

	NewHookType            EventType = "NEW_HOOK_TYPE"
	ChangeHookStatusType   EventType = "CHANGE_HOOK_STATUS_TYPE"
	ChangeHookEndpointType EventType = "CHANGE_HOOK_ENDPOINT_TYPE"
	ChangeHookSecretType   EventType = "CHANGE_HOOK_SECRET_TYPE"
	ChangeHookRetryType    EventType = "CHANGE_HOOK_RETRY_TYPE"
)

type Event struct {
	Channel   Channel     `json:"channel"`
	ID        string      `json:"id"`
	Attribute string      `json:"attribute"`
	Value     interface{} `json:"value"`
}

func (e Event) GetId() string {
	return e.ID
}

func (e Event) GetChannel() Channel {
	return e.Channel
}

func (e Event) GetAttribute() string {
	return e.Attribute
}

func (e Event) GetValue() interface{} {
	return e.Value
}

func (e Event) FromPayload(payload string) (Event, error) {
	var event Event
	err := json.Unmarshal([]byte(payload), &event)
	return event, err

}
func (e Event) ToPayload() (string, error) {
	data, err := json.Marshal(e)
	return string(data), err
}

func EventFromType(t EventType, attempt *Attempt, hook *Hook) (Event, error) {
	switch t {
	case NewWaitingAttemptType:
		if attempt == nil {
			return Event{}, errors.New(fmt.Sprintf("BuildEvent: Missing %s for %s type", "attempt", t))
		}
		return NewEvent(AttemptChannel, attempt.ID, "id", attempt.ID), nil
	case FlushWaitingAttemptType:
		if attempt == nil {
			return Event{}, errors.New(fmt.Sprintf("BuildEvent: Missing %s for %s type", "attempt", t))
		}
		return NewEvent(AttemptChannel, attempt.ID, "nexttry", "one"), nil
	case FlushWaitingAttemptsType:
		return NewEvent(AttemptChannel, "", "nexttry", "all"), nil
	case AbortWaitingAttemptType:
		if attempt == nil {
			return Event{}, errors.New(fmt.Sprintf("BuildEvent: Missing %s for %s type", "attempt", t))
		}
		return NewEvent(AttemptChannel, attempt.ID, "abort", "one"), nil

	case NewHookType:
		if hook == nil {
			return Event{}, errors.New(fmt.Sprintf("BuildEvent: Missing %s for %s type", "hook", t))
		}
		return NewEvent(HookChannel, hook.ID, "id", hook.ID), nil
	case ChangeHookStatusType:
		if hook == nil {
			return Event{}, errors.New(fmt.Sprintf("BuildEvent: Missing %s for %s type", "hook", t))
		}
		return NewEvent(HookChannel, hook.ID, "status", hook.Status), nil
	case ChangeHookEndpointType:
		if hook == nil {
			return Event{}, errors.New(fmt.Sprintf("BuildEvent: Missing %s for %s type", "hook", t))
		}
		return NewEvent(HookChannel, hook.ID, "endpoint", hook.Endpoint), nil
	case ChangeHookSecretType:
		if hook == nil {
			return Event{}, errors.New(fmt.Sprintf("BuildEvent: Missing %s for %s type", "hook", t))
		}
		return NewEvent(HookChannel, hook.ID, "secret", hook.Secret), nil
	case ChangeHookRetryType:
		if hook == nil {
			return Event{}, errors.New(fmt.Sprintf("BuildEvent: Missing %s for %s type", "hook", t))
		}
		return NewEvent(HookChannel, hook.ID, "retry", hook.Retry), nil
	}

	return Event{}, nil
}

func TypeFromEvent(ev Event) EventType {
	if ev.Channel == AttemptChannel {
		if ev.Attribute == "id" {
			return NewWaitingAttemptType
		}
		if ev.Attribute == "nexttry" {
			if ev.ID != "" {
				return FlushWaitingAttemptType
			}
			return FlushWaitingAttemptsType
		}
		if ev.Attribute == "abort" {
			if ev.ID != "" {
				return AbortWaitingAttemptType
			}
		}
	}

	if ev.Channel == HookChannel {
		if ev.Attribute == "id" {
			return NewHookType
		}
		if ev.Attribute == "status" {
			return ChangeHookStatusType
		}
		if ev.Attribute == "endpoint" {
			return ChangeHookEndpointType
		}
		if ev.Attribute == "secret" {
			return ChangeHookSecretType
		}
		if ev.Attribute == "retry" {
			return ChangeHookRetryType
		}

	}

	return ""
}

func NewEvent(channel Channel, id string, attribute string, value interface{}) Event {
	return Event{
		Channel:   channel,
		ID:        id,
		Attribute: attribute,
		Value:     value,
	}
}
