package main

import (
	"github.com/SashaCrofter/cjdngo"
	. "github.com/SashaCrofter/cjdngo/admin"
	"github.com/SashaCrofter/network"
	"math"
)

// Filter iterates through every route in the given table, and
// discards any Routes that are more than maxHops (discounted if <1)
// away from the origin, or that pass through any of the IP addresses
// as identified by endpoint. Any domain names that are intended to be
// used should be preresolved. It returns the remaining routes.
func Filter(table []*Route, maxHops int, endpoint []string) (filtered []*Route) {
	useHops := maxHops > 0

	if !useHops && len(endpoint) == 0 {
		// If we can avoid filtering altogether, by all means do.
		l.Println("No filters applied")
		return table
	}
	// Otherwise, report the arguments.
	l.Println("Filtering by hops:", useHops)
	l.Println("Number of end IPs:", len(endpoint))

	// Call it bad practice, but it's convenient, here.
	containsBadHop := func(hops []*Route) bool {
		if len(endpoint) < 1 || len(hops) < 1 {
			return false
		}
		// For every given hop
		for _, route := range hops {
			// check every string in endpoint.
			for _, endHop := range endpoint {
				if route.IP == endHop {
					// Only return true if at least one hop is
					// considered an end.
					l.Println("Discarding route that passes through",
						route.IP)
					return true
				}
			}
		}
		return false
	}

	for _, route := range table {
		hops := getNodesOnPath(table, route.Path)
		if !(useHops && len(hops) > maxHops) && !containsBadHop(hops) {
			// If there were few enough hops, and it did not pass
			// through any blacklisted nodes, then add it to the array
			// to be returned.
			filtered = append(filtered, route)
		}
		// Otherwise, it is discarded.
	}
	return
}

// ToNetwork converts from a routing table to Network form,
// establishing connections and determining routes. It invokes
// truncate() on the routing table in order to shorten the IPs.
func ToNetwork(routes []*Route) (n *network.Network) {
	n = &network.Network{
		Nodes: make(map[string]*network.Node),
	}
	truncate(routes)

	// For every route in the given table, then use FilterRoutes() to
	// get all of the hops on that route. Use the final hop in that
	// result as a peer. Then, filter duplicates.
	for _, route := range routes {
		node, isPresent := n.Nodes[route.IP]
		if !isPresent {
			// If the node isn't already present, make it.
			node = &network.Node{
				Connected:  make([]*network.Connection, 0),
				Attributes: make(map[string]interface{}),
			}
		}

		// Get all of the hops in that route
		hops := getNodesOnPath(routes, route.Path)
		// Because of self-routing bugs, it's necessary to look
		// through the list of hops.
		i := len(hops) - 1
		for {
			if i >= 0 {
				// Then, if there were any, (this should always be
				// true, but let's avoid runtime errors), append its
				// IP to the Connected attribute of the current Node.
				lastHop := hops[i].IP
				if route.IP != lastHop {
					l.Println("Node", route.IP, "connected to", lastHop)
					node.Connected = append(node.Connected,
						&network.Connection{
							Target: lastHop,
						})
					break
				} else {
					l.Println("Got self-connection on", route.IP)
				}
				i--
			} else {
				// If we get a zero-length path, report it, but don't
				// do anything.
				l.Println("Got broken path for", route.IP)
				break
			}
		}
		n.Nodes[route.IP] = node
	}

	// Now, we have to filter out duplicate connections.
	for _, node := range n.Nodes {
		// Copy all Connected entries into a map[string]interface{},
		// then copy them back.
		connectedMap := make(map[string]interface{})
		for _, connection := range node.Connected {
			connectedMap[connection.Target] = nil
		}
		node.Connected = make([]*network.Connection, 0, len(connectedMap))
		for k := range connectedMap {
			node.Connected = append(node.Connected,
				&network.Connection{
					Target: k,
				})
		}
	}
	return
}

func getNodesOnPath(table []*Route, target uint64) (hops []*Route) {
	for _, route := range table {

		g := 64 - uint64(math.Log2(float64(route.Path)))
		h := uint64(uint64(0xffffffffffffffff) >> g)

		if h&target == h&route.Path {
			hops = append(hops, route)
		}
	}
	return
}

// Convert all IPv6 addresses in a table into their shortest usable
// form. This just serves as a wrapper for cjdngo.Truncate().
func truncate(table []*Route) {
	ipMap := make(map[string]interface{})
	for _, route := range table {
		ipMap[route.IP] = nil
	}

	// Convert the map into an array.
	ips := make([]string, len(ipMap))
	var i int
	for ip := range ipMap {
		ips[i] = ip
		i++
	}

	t := cjdngo.Truncate(ips)
	for _, route := range table {
		route.IP = t[route.IP]
	}
}
