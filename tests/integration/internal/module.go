package internal

type Module struct {
	name           string
	createDatabase bool
	services       []service
}

func (m *Module) WithCreateDatabase() *Module {
	m.createDatabase = true
	return m
}

func (m *Module) WithServices(services ...service) *Module {
	for _, service := range services {
		m.services = append(m.services, service)
	}
	return m
}

func NewModule(name string) *Module {
	return &Module{
		name: name,
	}
}
