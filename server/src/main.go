package main

import (
	"fmt";
	"flag";
)

var settings = flag.String("c", "../settings.conf", "path to the configuration file")
var debug = flag.Bool("d", false, "set the debug mode of the program")

func main() {

	flag.Parse()
	fmt.Printf("setting path: %s\n", *settings)
	fmt.Printf("debug mode: %t\n", *debug)
}


