package main

import (
	//"encoding/json"
	"github.com/SashaCrofter/cjdngo"
	"github.com/SashaCrofter/cjdngo/admin"
	//"github.com/SashaCrofter/network-markup/nmparser"
	"log"
)

var (
	confLocation = []string{
		"/etc/cjdroute.conf",
		"/opt/cjdns/cjdroute.conf"} // Common config locations
)

func main() {
	// Get the configuration file
	var conf *cjdngo.Conf
	var err error
	for _, l := range confLocation {
		conf, err = cjdngo.ReadConf(l)
		if conf != nil && err == nil {
			break
		}
	}
	if conf == nil {
		log.Fatal("Configuration could not be loaded:", err)
	}

	// Connect to the admin interface
	cjdns, err := admin.Connect("", "", conf.Admin.Password)
	defer cjdns.Close() // Make sure that the connection is closed
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the routing table
	table := cjdns.DumpTable(-1)
	if len(table) == 0 {
		log.Fatalln("Could not dump the routing table.")
	}

	// Filter the routing table to contain only the relevant portions
	// of the table, then convert it to network format.
	_ = ToNetwork(Filter(table, 0, nil))
}
