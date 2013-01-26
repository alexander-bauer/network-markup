# Network Markup

*(Note: Network markup is in its early stages, and undergoing frequent
changes in specification and implementation, particularly of the
graphing style. Examples below are subject to improvement!)*

The basic goal of network markup is to allow the visualization (and
portable representation) of complex networks. My reason for developing
this particular form of markup is [necessity][meshnet]; I am involved in
peer-to-peer mesh networks, and so need to be able to quickly and
*easily* create mockups and "goal" visualizations.

For example, consider this paragraph. Alice is connected to Bob and
Steve. Bob is connected to Eve. Eve is connected to Steve and Alice.

This is not a particularly complex network, but it is, nonetheless, a
network. When converted to an easily-machine-readable JSON format, and
then made into an image using the [igraph][] library (via Python,)
it looks like this.

![simple example network][network demo]

This is generated *completely automatically*. The only part of this
demonstation that was not automatic was the actual network description,
which is clear and human-readable. See below for the [technical
details](#technical-details).

Simple, plain English markup languages have a great deal of potential
because of their extensibility. To be able to visualize and model both
complex and simple networks in a format that requires almost no overhead
is something that many people, programmers or not, might wish for.

## How to Use

Network markup is in its early stages. For example, a technical
specification doesn't yet exist, though it will be included in this
repository when written. The parser and grapher implementations are also
in development, but they have become stable enough to use.

### Installation

In order to write network markup, using this implementation of it, you
will first need to install the [Golang][] program included in this
repository. Once Go is properly installed, you will be able to use `go
install github.com/SashaCrofter/network-markup` to do this.

You will also need to have `python-igraph` installed. Details for doing
that can be found [here][igraph install].

### Composition

As of this writing, network markup sentences have three primary
components, meant to look like plain English. They are the subject,
attributes, and connections. Consider the sentence: "Alice is disabled
and connected to Bob."

This has three parts.

- Subject "*Alice*"
- Attributes "is *disabled*"
- Connections "and connected to *Bob*."

Network markup uses keywords, such as "is," and "connected" to determine
the intention of the sentence. Because of this, the above sentence is
equivalent to "Alice is connected to Bob and is disabled." It is not,
however, equivalent to "Alice is connected to Bob and disabled." My
implementation would interpret the latter sentence and create a graph in
which Alice is connected to the nodes "Bob" and "disabled."

[Examples][] are included in this repository in both `.nm` and
`.json` form. The specification and implementation may change in their
details, but they should remain mostly the same.

## Technical Details

Even the JSON format is relatively trivial, though. This was the output
of the markup parser ([nmparser][]), with indentation enabled.

```json
{
	"Alice": {
		"connected": [
			"Bob",
			"Steve"
		]
	},
	"Bob": {
		"connected": [
			"Eve"
		]
	},
	"Eve": {
		"connected": [
			"Steve",
			"Alice"
		]
	},
	"Steve": {
		"connected": []
	}
}
```

This is an open format which can be created by any program, automated or
not. The markup parser itself, in the implementation in this repository,
is written in [Golang][], with the back-end library exposed for easy
external usage. ([Linked previously][nmparser])

The python script which is used to convert the above JSON to the
graphical representation is [also in this repository][python script],
and also simply an implementation. The library which is used,
[igraph][], is in C, and has both R and Python interfaces, the
latter of which I am using.


[network demo]: http://i.imgur.com/e5hoiqC.png
	"Simple network demo"

[nmparser]: http://godoc.org/github.com/SashaCrofter/network-markup/nmparser
	"Network Markup Parser Go Library"

[python script]: python/jsonToIGraph.py
[examples]: example/

[meshnet]: https://maryland.projectmeshnet.org
	"Maryland Mesh Homepage"

[igraph]: http://igraph.sourceforge.net
	"igraph Homepage"
[igraph install]: http://igraph.wikidot.com/installing-python-igraph-on-linux
	"igraph Installation"

[golang]: http://golang.org
	"Golang Homepage"