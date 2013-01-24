#!/usr/bin/env python
import igraph
import json
import sys

if(len(sys.argv) == 1):
    print('You need to supply a JSON file.')
    exit(1)

# Get the contents of the relevant file
j = json.load(open(sys.argv[1], 'r'))

# Initialize the graph with n = nodes in the input file
g = igraph.Graph(0, directed=True)

# Add all of the vertices, so that there are no unknown objects. They
# are named according to their name in the JSON
i = 0
for name in j:
    print "Adding vertex: " + name
    g.add_vertex(name)

    # Now, g.vs[i] refers to the most recent vertex
    g.vs[i]["label"] = name
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

# Finally, plot the graph
fig = igraph.Plot(bbox=(480, 480), background="white")
fig.add(g, layout="fr")
fig.show()
