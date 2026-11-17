package main

import (
	"flag"
	"fmt"
	"os"
)

// parseFlags handles the small set of operational flags the binary supports.
func parseFlags(args []string) (printVersion bool, err error) {
	fs := flag.NewFlagSet("ledger", flag.ContinueOnError)
	fs.BoolVar(&printVersion, "version", false, "print build information and exit")
	fs.SetOutput(os.Stderr)
	if err = fs.Parse(args); err != nil {
		return false, err
	}
	if printVersion {
		info := buildInfo()
		fmt.Printf("ledger %s (%s)\n", info["version"], info["commit"])
	}
	return printVersion, nil
}
