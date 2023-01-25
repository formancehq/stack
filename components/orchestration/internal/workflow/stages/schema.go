package stages

var stages = map[string]Stage{}

type Stage interface {
	GetWorkflow() any
}

func Register(name string, stage Stage) {
	stages[name] = stage
}

func All() map[string]Stage {
	return stages
}

func Get(name string) Stage {
	return stages[name]
}
