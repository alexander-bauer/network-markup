package main

import (
	"encoding/json"
	"flag"
	"github.com/SashaCrofter/cjdngo"
	"github.com/SashaCrofter/cjdngo/admin"
	//	"github.com/SashaCrofter/network-markup/nmparser"
	"io/ioutil"
	"log"
	"os"
)

var (
	confLocation = []string{
		"/etc/cjdroute.conf",
		"/opt/cjdns/cjdroute.conf"} // Common config locations

	l      *log.Logger // System wide logger
	output = os.Stdout // Output system
)

// Flags
var (
	fLog = flag.Bool("l", false, "Enable logging output")
)

func main() {
	flag.Parse()
	if *fLog {
		l = log.New(os.Stdout, "", log.Ltime)
	} else {
		l = log.New(ioutil.Discard, "", 0)
	}

	// Use the arguments as an array of blacklisted nodes, for
	// filtering. This will probably be changed in the future.
	blacklist := flag.Args() // Does not include argv[0]

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
		log.Fatalln("Configuration could not be loaded:", err)
	}

	// Connect to the admin interface
	cjdns, err := admin.Connect("", "", conf.Admin.Password)
	defer cjdns.Close() // Make sure that the connection is closed
	if err != nil {
		log.Fatalln(err)
	}

	// Retrieve the routing table
	table := cjdns.DumpTable(-1)
	if len(table) == 0 {
		log.Fatalln("Could not dump the routing table.")
	}

	// Filter the routing table to contain only the relevant portions
	// of the table, then convert it to network format.
	network := ToNetwork(Filter(table, 0, blacklist))
	err = json.NewEncoder(output).Encode(network)
	if err != nil {
		log.Fatalln("Encountered error while marshalling output:",
			err)
	}
}
