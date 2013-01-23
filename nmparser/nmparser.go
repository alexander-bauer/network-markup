package nmparser

import (
	"strings"
)

type Network map[string]*Node

type Node struct {
	Connected  []string
	Attributes map[string]interface{}
}

func Parse(nm string) (n *Network) {
	tokenize(strings.NewReader(nm))
	return nil
}
