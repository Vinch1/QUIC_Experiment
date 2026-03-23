package cluster

type Node struct {
	ID      string
	Address string
}

type Topology struct {
	Nodes []Node
}

func (t Topology) Size() int {
	return len(t.Nodes)
}
