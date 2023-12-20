package messages

type Messages struct {
	stackURL string
}

func NewMessages(stackURL string) *Messages {
	return &Messages{
		stackURL: stackURL,
	}
}
