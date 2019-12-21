package main 

import (
	"fmt"
	"net"
	"strings"
)

const (
    CONN_HOST = "localhost"
)

func runServer(cfg Config) {
	host := strings.Join([]string {CONN_HOST, cfg.Port}, ":")
	if *debug {
		fmt.Printf("Host: %s\n", host)
	}
	conn, err := net.Listen("tcp", host)
	if err != nil {
		fatal(fmt.Sprintf("Can't listen to port %d at localhost\n", cfg.Port))
	}
	defer conn.Close()

	fmt.Printf("Starting server...\n")
	fmt.Printf("Listening on Port %s\n", cfg.Port)
}