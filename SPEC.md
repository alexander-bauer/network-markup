# Network Markup

Network markup is meant to be an open standard for representing
complex networks in plain English. The below specification is based
off of the existing implementation, available at its
[repository][]. **It is expected to change as the the implementation
advances, then stabilize at a later date.**

For expansion on the premise, please see the [readme][]. Examples are
available [here][examples].

## Parsing

### Components

There are three main components of each declaration. There are the
node that is being declared, the attributes of that node, and the
nodes to which that node is connected. These can be seen clearly in
the following sentence.

	Alice is disabled and is connected to Bob.

The node Alice is declared here, has the attribute "disabled," and is
connected to another node, Bob. A parser should read it as follows.

	Alice is disabled and is connected to Bob.
	^     ^  ^        ^      ^            ^
	|     |  |        |      |            Connection
	|     |  |        |      Connections separator
    |     |  |        Ignored separator
	|     |  Attribute
	|     Attributes separator
	Declaration

Words such as "and" and "is," and phrases like "is connected to," must
be used by the parser for context. This is not only so that network
markup can resemble normal English, but also so that sentences do not
*need* to follow an artificial structure. For example, the following
two sentences should be exactly equivalent.

	Alice is disabled and is connected to Bob.
	Alice is connected to Bob and is disabled.

They must also be used for context to support multi-word
identifiers. A network markup parser should be able to understand the
following.

	Data Center One is connected to Bob, Eve, and Data Center Two.

All multi-word identifiers should be understood correctly, including
attributes.

    Data Center Two is high bandwidth.

### Implicit Declaration

In many cases, nodes are not declared explicitly (*ie*. not in their
own sentences). This means that they are declared implicitly by
mentions after other nodes. Consider the following.

	Alice is connected to Bob.

Alice is declared explicitly, but Bob is not. Therefore, the parser
has to make the assumption that Bob is another node, and include him
in the output.

## Output Format

The preferred network markup output format is currently unnamed. It
would normally be outside of the scope of this specification, but
because of its simplicity, it will be covered briefly.

The format is simply a particular structure of a JSON (JavaScript
Object Notation) object, which maps node names to connections and
attributes. The template format is this.

	{
		"<node one>": {
			"connected": [ "<node two>", "<node three>" ]
		},
		"<node two>": {
			"connected": [],
			"attributes": {
				"<attribute name>": <value of any type>
			}
		},
		"<node three>" {
			"connected": []
		}
	}

This is, in Python terms, a dictionary which maps strings to other
dictionaries, which map strings to either arrays or dictionaries.

The node names, such as "Alice" or "Bob," are required, and there
*cannot* be implicit declarations.

The `connected` item is also required, but it doesn't need to contain
any items. Note that it is an array type.

The `attributes` item is *not* required, but this **may change in the
future**. If included, it is a dictionary which maps strings to values
of any type.

## Conclusion

This specification is in early development. It is subject to change,
and very much subject to improvement. Please submit requests or
clarifications (regarding any part of this project, including the
specification) to the [GitHub issues page][issues].

[repository]: https://github.com/SashaCrofter/network-markup
[readme]: https://github.com/SashaCrofter/network-markup#readme
[examples]: https://github.com/SashaCrofter/network-markup/tree/development/example
[issues]: https://github.com/SashaCrofter/network-markup/issues
