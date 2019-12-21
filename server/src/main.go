package main

import (
	"fmt";
	"flag";
)

var configPath = flag.String("c", "settings.conf", "path to the configuration file")
var debug = flag.Bool("d", true, "set the debug mode of the program")

func main() {
	flag.Parse()
	config := readConfigFile(*configPath)
	if *debug {
		fmt.Printf("Setting path: %s\n", *configPath)
		fmt.Printf("Port: %s\n", config.Port)
		fmt.Printf("Debug mode: %t\n", *debug)
	}
	
	runServer(config)


}


