package dockerutils

type ContainerCmd struct {
	ID      string // the container name or ID
	Command string
	Args    []string
	Env     []string
}
