package main

import (
	. "github.com/SashaCrofter/cjdngo/admin"
	"github.com/SashaCrofter/network-markup/nmparser"
	"math"
)

// Filter iterates through every route in the given table, and
// discards any Routes that are more than maxHops (discounted if <1)
// away from the origin, or that pass through any of the IP addresses
// as identified by doNotGoThrough. Any domain names that are intended
// to be used should be preresolved. It returns the remaining routes.
func Filter(table []*Route, maxHops int, doNotGoThrough []string) (filtered []*Route) {
	useHops := maxHops > 0

	if !useHops && len(doNotGoThrough) == 0 {
		// If we can avoid filtering altogether, by all means do.
		l.Println("No filters applied")
		return table
	}
	// Otherwise, report the arguments.
	l.Println("Filtering by hops:", useHops)
	l.Println("Number of bad IPs:", len(doNotGoThrough))

	// Call it bad practice, but it's convenient, here.
	containsBadHop := func(hops []*Route) bool {
		if len(doNotGoThrough) < 1 {
			return false
		}
		// For every given hop,
		for _, route := range hops {
			// check every string in doNotGoThrough.
			for _, badHop := range doNotGoThrough {
				if route.IP == badHop {
					// Only return true if at least one hop is
					// considered bad.
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

func ToNetwork(routes []*Route) (network nmparser.Network) {
	network = make(nmparser.Network)
	// For every route in the given table, then use FilterRoutes() to
	// get all of the hops on that route. Use the final hop in that
	// result as a peer. Then, filter duplicates.
	for _, route := range routes {
		node, isPresent := network[route.IP]
		if !isPresent {
			// If the node isn't already present, make it.
			node = &nmparser.Node{
				Connected:  make([]string, 0),
				Attributes: make(map[string]interface{}),
			}
		}

		// Get all of the hops in that route
		hops := getNodesOnPath(routes, route.Path)
		if len(hops) > 0 {
			// Then, if there were any, (this should always be true,
			// but let's avoid runtime errors), append its IP to the
			// Connected attribute of the current Node.
			lastHop := hops[len(hops)-1].IP
			if route.IP != lastHop {
				l.Println("Node", route.IP, "connected to", lastHop)
				node.Connected = append(node.Connected, lastHop)
			} else {
				l.Println("Got self-connection on", route.IP)
			}
		} else {
			// If we get a zero-length path, report it, but don't do
			// anything.
			l.Println("Got zero-hop path for", route.IP)
		}
		network[route.IP] = node
	}
	/*
		// Now, we have to filter out duplicate connections.
		for _, node := range network {
			// Copy all Connected entries into a map[string]interface{},
			// then copy them back.
			connectedMap := make(map[string]interface{})
			for _, connection := range node.Connected {
				connectedMap[connection] = nil
			}
			node.Connected = make([]string, 0, len(connectedMap))
			for k := range connectedMap {
				node.Connected = append(node.Connected, k)
			}
		}*/
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
