package ingester

type StateLabel string

const (
	StateLabelInit  StateLabel = "INIT"
	StateLabelReady StateLabel = "READY"
	StateLabelPause StateLabel = "PAUSE"
	StateLabelStop  StateLabel = "STOP"
)

type State struct {
	Label StateLabel `json:"label"`
	// Cursor can be set only if Label == StateLabelInit
	Cursor string `json:"cursor,omitempty"`
	// PreviousState can be set only if Label == StateLabelPause or Label == StateLabelStop
	PreviousState *State `json:"previousState,omitempty"`
}

func (s State) String() string {
	switch s.Label {
	case StateLabelInit:
		return "INIT"
	case StateLabelReady:
		return "READY"
	case StateLabelPause:
		return "PAUSE"
	case StateLabelStop:
		return "STOP"
	default:
		return "UNKNOWN_STATE"
	}
}

func NewInitStateWithCursor(cursor string) State {
	return State{
		Label:  StateLabelInit,
		Cursor: cursor,
	}
}

func NewInitState() State {
	return NewInitStateWithCursor("")
}

func NewStopState(previousState State) State {
	return State{
		Label:         StateLabelStop,
		PreviousState: &previousState,
	}
}

func NewReadyState() State {
	return State{
		Label: StateLabelReady,
	}
}

func NewPauseState(previousState State) State {
	return State{
		Label:         StateLabelPause,
		PreviousState: &previousState,
	}
}
