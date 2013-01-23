#!/usr/bin/Rscript

jsonToNM <- function(jfile, directed=TRUE) {
  # Requires "rjson" and "igraph"
  
  nm.json <- fromJSON(file=jfile)
  nm.graph <- c()
  
  # Expand the $network row into igraph form
  for(i in 1:length(nm.json$network)) {
    if(length(nm.json$network[[i]] > 0)) {
      for(j in 1:length(nm.json$network[[i]])) {
        nm.graph <- c(nm.graph, i, nm.json$network[[i]][j])
      }
    }
  }
  
  # Create the network map object
  g <- graph(nm.graph, directed=directed)
  # Apply the names
  V(g)$name <- nm.json$names
  V(g)$label <- V(g)$name
  
  plot(g, vertex.label.dist=1.5)
}
