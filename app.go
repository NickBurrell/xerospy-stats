package main

import (
	"flag"
	"github.com/zero-frost/xerospy-stats/app"
)

func main() {

	var debugFlag bool

	// Parse command line flags
	flag.BoolVar(&debugFlag, "debug", false, "Enables debug print-outs")
	flag.Parse()

	// Initialize server and templating engine
	server := app.Server{DebugMode: debugFlag}
	server.InitServer()
	server.RunServer()
}
