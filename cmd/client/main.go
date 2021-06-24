// Generates binary for the client

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vadimegorov13/reminder-go/client"
)

var (
	backendURIFlag = flag.String("backend", "http://localhost:5000", "Backend API URL")
	helpFlag       = flag.Bool("help", false, "Display all the commands and their usage")
)

func main() {
	flag.Parse()
	s := client.NewSwitch(*backendURIFlag)

	if *helpFlag || len(os.Args) == 1 {
		s.Help()
		return
	}

	err := s.Switch()
	if err != nil {
		fmt.Printf("cmd switch error: %s", err)
		os.Exit(2)
	}
}
