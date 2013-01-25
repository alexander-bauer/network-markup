package nmparser

import (
	"errors"
	"fmt"
	"strings"
)

type Network map[string]*Node

type Node struct {
	Connected []string `json:"connected"` // All of
	// the connected nodes
	Attributes map[string]interface{} `json:"attributes,omitempty"` // Any
	// attributes (unused)
}

var (
	MultipleSubjectError = errors.New(`Multiple nodes were declared in
the same sentence: %s`)
)

// Parts of the sentences
const (
	partDeclaration = iota // Pre-"is"
	partAttributes         // Pre-"is connected to" and post-"is"
	partConnections        // Post-"is connected to"
)

func Parse(nm string) (n Network, err error) {
	tokens := tokenize(strings.NewReader(nm))

	n = make(Network)

	var subject *Node // Node declared in this sentence
	var part int      // Current part of the sentence
	var negated bool  // Whether the part of speech is negated
	for _, t := range tokens {
		// If we've hit a newline or other separator, though, we need
		// to mark ourselves as at the beginning of the sentence.
		if t.Id == tokenNEWLINE {
			subject = nil
			part = partDeclaration
		} else if t.Id == tokenIS {
			// This is true, no matter what the previous part
			// was. Consider "Alice is connected to Bob, and is
			// disabled."
			part = partAttributes
		} else if t.Id == tokenCONNTO {
			// Likewise, this will always work.
			part = partConnections
		} else if t.Id == tokenSEPARATOR {
			// The separators are simply grammatical constructs, and
			// don't mean much to us.
		} else if t.Id == tokenNEGATOR {
			// Negate the following attribute. This is automatically
			// reset at the end of this block, so we set true here and
			// jump to the start of the loop.
			negated = true
			continue
		} else if part == partDeclaration && t.Id == tokenIDENT {
			// If we're in the early part of the sentence, we should be
			// looking for the newly declared node. "Alice is..."
			if subject != nil {
				// If we reach this part, and subject isn't nil,
				// there's something very wrong.
				return nil, fmt.Errorf(MultipleSubjectError.Error(),
					t.Literal)
			}

			// We must declare node to avoid an additional map
			// lookup.
			var node *Node
			if _, isPresent := n[t.Literal]; !isPresent {
				// If the node hasn't already been declared, create it
				// here. Otherwise, don't.
				node = &Node{
					Connected:  make([]string, 0),
					Attributes: make(map[string]interface{}),
				}
				n[t.Literal] = node
			} else {
				node = n[t.Literal]
			}
			// Set subject to our pointer, for easy access.
			subject = node
		} else if part == partAttributes && t.Id == tokenIDENT {
			// Dumbly mark the attribute name, and "true" if it is not
			// negated. Consider "Alice is disabled" and "Alice is not
			// disabled." The former implies
			//
			//     "disabled": true
			//
			// and the latter implies
			//
			//     "disabled": false
			subject.Attributes[t.Literal] = !negated
		} else if part == partConnections && t.Id == tokenIDENT {
			// If we're after the "connected", then we need to add any
			// new nodes to the most recent Node's "connected" block.
			// Additionally, we will need to recognize any attributes.
			if t.Id == tokenIDENT {
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
		// Set negated to false here, so that one "not" doesn't negate
		// the whole structure.
		negated = false
	}
	return
}
