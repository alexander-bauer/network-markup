package nmparser

import (
	"strings"
)

type Network map[string]*Node

type Node struct {
	Connected  []string               // All of the connected nodes
	Attributes map[string]interface{} // Any attributes (unused)
}

func Parse(nm string) (n Network) {
	tokens := tokenize(strings.NewReader(nm))

	n = make(Network)

	var isAfterSeparator bool
	var mostRecentNode *Node
	for _, t := range tokens {
		// BUG(DuoNoxSol): This means of identifying nodes
		// will break when confronted with spaces.

		// The first step in the parsing is determine what part of the
		// sentence we're currently in. We're either before the "is
		// connected to" or after it. If we're in it, we just need to
		// set isAfterSeparator real quick.
		if t.Id == tokenCONNTO {
			isAfterSeparator = true
			continue
		} else if t.Id == tokenNEWLINE {
			// If we've hit a newline or other separator, though, we
			// need to mark ourselves as at the beginning of the
			// sentence.
			isAfterSeparator = false
			continue
		}

		// If we're before the separator, we should add any new nodes
		// we encounter.
		if !isAfterSeparator {
			if t.Id == tokenNODE {
				// Add the new node and its identifier in the Network
				mostRecentNode = &Node{
					Connected:  make([]string, 0),
					Attributes: make(map[string]interface{}),
				}
				n[t.Literal] = mostRecentNode
			}
			// If the token is anything else, we can discard it for
			// the time being.
		} else {
			// If we're after the separator, then we need to add any
			// new nodes to the most recent Node's "connected" block.
			if t.Id == tokenNODE {
				mostRecentNode.Connected = append(
					mostRecentNode.Connected, t.Literal)
			}
		}
	}
	return
}
