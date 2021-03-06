package nmparser

import (
	"errors"
	"fmt"
	"github.com/SashaCrofter/network"
	"strings"
)

var (
	ParserError          = errors.New("Parsing error")
	MultipleSubjectError = errors.New(`Multiple nodes were declared in
the same sentence: %s`)
)

// Parts of the sentences
const (
	partDeclaration        = iota // Pre-"is"
	partAttributes                // Pre-"is connected to" and post-"is"
	partConnections               // Post-"is connected to"
	partConnectionModifier        // The latter part of "is connected to
	// ... by ..."
)

func Parse(nm string) (n *network.Network, err error) {
	tokens := tokenize(strings.NewReader(nm))

	n = &network.Network{
		Nodes: make(map[string]*network.Node),
	}

	var subject *network.Node      // Node declared in this sentence
	var object *network.Connection // Node that the subject is connected
	// to
	var part int     // Current part of the sentence
	var negated bool // Whether the part of speech is negated
	for _, t := range tokens {
		// If we've hit a newline or other separator, though, we need to
		// mark ourselves as at the beginning of the sentence.
		if t.Id == tokenNEWLINE {
			subject, object = nil, nil
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
			var node *network.Node
			if _, isPresent := n.Nodes[t.Literal]; !isPresent {
				// If the node hasn't already been declared, create it
				// here. Otherwise, don't.
				node = &network.Node{
					Connected:  make([]*network.Connection, 0),
					Attributes: make(map[string]interface{}),
				}
				n.Nodes[t.Literal] = node
			} else {
				node = n.Nodes[t.Literal]
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
				object = &network.Connection{
					Target: t.Literal,
				}
				subject.Connected = append(subject.Connected, object)
				// Now that it's been connected, we need to make note
				// of the node if it doesn't already exist.
				if _, isPresent := n.Nodes[t.Literal]; !isPresent {
					n.Nodes[t.Literal] = &network.Node{
						Connected:  make([]*network.Connection, 0),
						Attributes: make(map[string]interface{}),
					}
				}
			}
		} else if part == partConnections && t.Id == tokenMOD {
			// We have presumably reached the beginning of a connection
			// modifier section, which will terminate after the next
			// tokenIDENT.
			part = partConnectionModifier
		} else if part == partConnectionModifier && t.Id == tokenIDENT {
			// TODO(DuoNoxSol): This will not support multiple
			// connection attributes.

			// If we've reached the identifier for the "by" section, add
			// it to the connection's Attributes field. If the node or
			// the connection hasn't been declared, though, something is
			// very wrong.
			if subject == nil || object == nil {
				return nil, ParserError
			}

			// If the attributes map doesn't exist, make it.
			if object.Attributes == nil {
				object.Attributes = make(map[string]interface{})
			}

			// Set the "label" attribute
			object.Attributes["label"] = t.Literal
			// and finally set the part back to partConnections, so that
			// multiple connections can be declared, even if they have
			// attributes.
			part = partConnections
		}
		// Set negated to false here, so that one "not" doesn't negate
		// the whole structure.
		negated = false
	}
	return
}
