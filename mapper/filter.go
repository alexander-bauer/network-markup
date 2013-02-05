package main

import (
	. "github.com/SashaCrofter/cjdngo/admin"
	"github.com/SashaCrofter/network-markup/nmparser"
)

// Filter iterates through every route in the given table, and
// discards any Routes that are more than maxHops (discounted if <1)
// away from the origin, or that pass through any of the IP addresses
// as identified by doNotGoThrough. Any domain names that are intended
// to be used should be preresolved. It returns the remaining routes.
func Filter(table []*Route, maxHops int, doNotGoThrough []string) (filtered []*Route) {
	// Once all routes have been retrieved, and the irrelevant ones
	// filtered, then you need to, in ToNetwork(), build a map of
	// connections (to filter out duplicates) and convert to an
	// array. Then all will be well.
	useHops := maxHops > 0

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
					return true
				}
			}
		}
		return false
	}

	for _, route := range table {
		hops := FilterRoutes(table, route.IP, 0, 0)
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

func ToNetwork(routes []*Route) (network *nmparser.Network) {
	println(len(routes))
	return nil
}
