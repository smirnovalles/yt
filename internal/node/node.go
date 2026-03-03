package node

import "fmt"

type Node struct {
	name string
}

func (n *Node) Print() string {
	return fmt.Sprintf("Node name is %s", n.name)
}

func NewNode() *Node {
	nn := Node{name: "node one"}
	return &nn
}
