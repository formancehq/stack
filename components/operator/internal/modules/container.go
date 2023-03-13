package modules

type Container struct {
	Command              []string
	Args                 []string
	Env                  ContainerEnv
	Image                string
	Name                 string
	Liveness             Liveness
	DisableRollingUpdate bool
}
