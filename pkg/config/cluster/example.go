package config

// A default cluster config

var (
	// DefaultClusterConfig :cluster config
	DefaultClusterConfig = Cluster{
		Name: DefaultClusterName ,
		Nodes: []Node {
			{Label: "server"},
			{Label: "worker"},
			{Label: "worker"},
		},
	}
)