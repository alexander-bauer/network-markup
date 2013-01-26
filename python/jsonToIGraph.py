#!/usr/bin/env python
import igraph
import json
import sys
import attributes

# Constants
bbox = (600, 600)
directedness = True
background = "white"

# Variables that need to be declared
output = ""

if(len(sys.argv) == 1):
    print('You need to supply a JSON file.')
    exit(1)

if len(sys.argv) > 2:
    # If a second argument is supplied, interpret it as an output
    # file.
    output = sys.argv[2]

# Get the contents of the relevant file
j = json.load(open(sys.argv[1], 'r'))

# Initialize the graph with n = nodes in the input file
g = igraph.Graph(0, directed=directedness)

# Add all of the vertices, so that there are no unknown objects. They
# are named according to their name in the JSON
i = 0
for name, node in j.items():
    
    print "Adding vertex: " + name
    g.add_vertex(name)

    # Now, g.vs[i] refers to the most recent vertex
    attributes.apply(g.vs[i], name, node)
    i += 1
    

# Add the edges by iterating through the "connected" field of the
# Node, and adding an edge in the graph for each
for name, node in j.items():
    # This needs to be wrapped in a try/except block, because it will
    # fail if node["connected"] doesn't exist
    try:
        for target in node["connected"]:
            print "Adding edge: " + name + "--" + target
            g.add_edge(name, target)
    except Exception:
        pass

# Make sure all edges are straight
g.es["curved"] = False

# Finally, plot the graph
if len(output) > 0:
    print "Writing to " + output
    igraph.plot(g, output, bbox=bbox, background=background, margin=50)
else:
    print "Displaying plot"
    igraph.plot(g, bbox=bbox, background=background, margin=50)
