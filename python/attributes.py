import igraph

# Applies known attributes, such as "disabled" to vertices by changing
# them into igraph-recognized values. The "label" attribute is
# automatically set to the supplied name. It returns the modified
# vertex.
def apply(vertex, name, node):
    vertex["label"] = name
    try:
        attributes = node["attributes"]
    except Exception:
        return vertex

    for name, key in attributes.items():
        # If the node is disabled, then grey it out.
        if name == 'disabled':
            if key == True:
                vertex["color"] = "grey"
    return vertex
