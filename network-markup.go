package main

import (
	"encoding/json"
	"flag"
	"github.com/SashaCrofter/network-markup/nmparser"
	"io/ioutil"
	"log"
	"os"
)

var (
	l *log.Logger // Standard logger

	fLog = flag.Bool("log", false, "enable logging")
)

func main() {
	flag.Parse()

	if *fLog {
		l = log.New(os.Stdout, "", 0)
	} else {
		l = log.New(ioutil.Discard, "", 0)
	}

	// If a filename isn't provided, then we need to fail here
	if flag.NArg() < 1 {
		log.Fatal("Not enough arguments. Please provide a filename.")
	}

	// Read the file's contents from the filesystem
	nmb, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	network := nmparser.Parse(string(nmb))
	jnet, err := json.MarshalIndent(network, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	println(string(jnet))
}
