package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var configPath = flag.String("c", "settings.conf", "path to the configuration file")
var debug = flag.Bool("d", true, "set the debug mode of the program")
var dataPresist = flag.Bool("p", true, "whether persist the data on filesystem")

func singnalHandling(db *DB) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		sig := <-sigs
		fmt.Printf("\n[%s]Persisting data to the file system...\n", sig)
		db.savedToFileSystem()
		fmt.Println("Data successfully saved to the file system")
		os.Exit(0)
	}()
}

func main() {
	flag.Parse()
	config := readConfigFile(*configPath)

	db := initDB(config)
	if *dataPresist {
		singnalHandling(db)
	}
	if *debug {
		fmt.Printf("Setting path: %s\n", *configPath)
		fmt.Printf("Port: %s\n", config.Port)
		fmt.Printf("Debug mode: %t\n", *debug)
	}
	runServer(config, db)
}
