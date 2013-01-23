#!/usr/bin/Rscript
require("rjson")
require("igraph")

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
