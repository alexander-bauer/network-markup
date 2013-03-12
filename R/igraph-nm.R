#!/usr/bin/env Rscript
require("rjson")
require("igraph")
require("getopt")

arguments <- rbind(c("file", "f", 1, "character"),
		  	 	   c("output", "o", 2, "character"),
				   c("directedness", "d", 2, "logical"))

args <- getopt(arguments)

background <- "white"

j <- fromJSON(file=args[["file"]])
g <- graph.empty(directed=args[["directedness"]])

# Adding vertices
print("Creating vertices")
names <- names(j$nodes)
for(n in 1:length(j$nodes)) {
  print(paste("Adding vertex:", names[n]))
  add.vertices(g, 1, attr=list(label=names[n]))
}

# Now add edges
print("Creating edges")
for(n in 1:length(j$nodes)) {
  print(paste("For node", n))
  print(j$nodes[n])
  print(j$nodes[n]$connected)
  for(connection in j$nodes[n]$connected) {
	print(paste("Adding edge:", names[n], "--", connection$target))
  	add.edge(g, edge(names[n], connection$target),
	  attr=list(label=connection$attributes$label))
  }
}

quit(0, save="no")

## Old Code

graphNJSON <- function(jfile, directed=TRUE) {
  # Requires "rjson" and "igraph"
  
  nm.json <- fromJSON(file=jfile)
  nm.graph <- c()
  
  # Initialize the graph with the given nodes
  g <- graph.empty(n=length(nm.json), directed=directed)
  # Add their names
  V(g)$name <- names(nm.json)
  V(g)$label <- V(g)$name
  
  # Now, add the edges
  for(i in 1:length(nm.json)) {
    # If the node has a "connected" field,
    # then we note the connections by looking
    # the names up.
    if(length(nm.json[[i]]$connected > 0)) {
      for(j in 1:length(nm.json[[i]]$connected)) {
        # Add the entry
        g <- g + edge(names(nm.json)[i],
                        nm.json[[i]]$connected[j])
      }
    }
  }
  
  plot(g, vertex.label.dist=1.5)
}
