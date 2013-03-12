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
g <- graph.empty(n=length(j$nodes), directed=args[["directedness"]])
# Add their names
my.names <- V(g)$label <- V(g)$name <- names(j$nodes)

g <- add.edges(g, c("Oberon", "Theseus"), attr=list(label="test"))

# Now add edges
for(n in 1:length(j$nodes)){
  node <- j$nodes[n][[1]]
  for(connection in node$connected){
    g <- add.edges(g, c(my.names[n], connection$target),
                      attr=connection$attributes)
  }
}

plot(g)
