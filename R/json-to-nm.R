#!/usr/bin/Rscript
# Check command line arguments
if(length(jfile) == 0) {
  jfile <- commandArgs(TRUE)[1]
}
if(length(jfile) == 0) {
    # Exit if the file is not supplied
    print("You must supply a file argument.")
    quit(1)
}

# Load required libraries
require("rjson")
require("igraph")

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
g <- graph(nm.graph)
# Apply the names
V(g)$name <- nm.json$names
V(g)$label <- V(g)$name

plot(g, vertex.label.dist=1.5)
