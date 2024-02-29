package workflow

type Input struct {
	Workflow  Workflow          `json:"workflow"`
	Variables map[string]string `json:"variables"`
}
