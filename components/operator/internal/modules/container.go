package modules

type Container struct {
	Command              []string
	Args                 []string
	Env                  ContainerEnv
	Image                string
	Name                 string
	DisableRollingUpdate bool
}
