package nmparser

import (
	"strings"
)

type Network map[string]*Node

type Node struct {
	Connected []string `json:"connected"` // All of
	// the connected nodes
	Attributes map[string]interface{} `json:"attributes,omitempty"` // Any
	// attributes (unused)
}

func Parse(nm string) (n Network) {
	tokens := tokenize(strings.NewReader(nm))

	n = make(Network)

	var subject *Node
	for _, t := range tokens {
		// BUG(DuoNoxSol): This means of identifying nodes
		// will break when confronted with spaces.

		// If we've hit a newline or other separator, though, we need
		// to mark ourselves as at the beginning of the sentence.
		if t.Id == tokenNEWLINE {
			subject = nil
			continue
		}

		// If we're before the separator, we should add any new nodes
		// we encounter.
		if subject == nil {
			if t.Id == tokenNODE {
				// Add the new node and its identifier in the Network
				var node *Node
				if _, isPresent := n[t.Literal]; !isPresent {
					node = &Node{
						Connected:  make([]string, 0),
						Attributes: make(map[string]interface{}),
					}
					n[t.Literal] = node
				} else {
					node = n[t.Literal]
				}
				subject = node
			}
			// If the token is anything else, we can discard it for
			// the time being.
		} else {
			// If we're after the separator, then we need to add any
			// new nodes to the most recent Node's "connected" block.
			// Additionally, we will need to recognize any
			// attributes.
			if t.Id == tokenNODE {
				subject.Connected = append(
					subject.Connected, t.Literal)
				// Now that it's been connected, we need to make note
				// of the node if it doesn't already exist.
				if _, isPresent := n[t.Literal]; !isPresent {
					n[t.Literal] = &Node{
						Connected:  make([]string, 0),
						Attributes: make(map[string]interface{}),
					}
				}
			}
		}
	}
	return
}
