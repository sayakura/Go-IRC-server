package irc

import (
	"fmt";
	"flag";
)

var settings = flag.String("c", "../settings.conf", "path to the configuration file")
var debug = flag.Bool("d", false, "set the debug mode of the program")

func main() {

	flag.Parse()
	fmt.print(settings)
	fmt.print(debug)
}


