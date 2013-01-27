import igraph

attb = {'disabled': {'color': 'grey'}}

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

    for aName, aVal in attributes.items():
        # If the attb dictionary contains the relevant attribute, then
        # apply all changes implied.
        try:
            changes = attb[aName]
            for cName, cVal in changes.items():
                vertex[cName] = cVal
        except Exception:
            print "Adding attribute: using literal " + str(aName)
            # If the change isn't known, though, try to apply it
            # anyway. This is for cases in which the attributes are
            # literal, like {'color':'blue'}
            vertex[str(aName)] = aVal
    return vertex
