package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	// Display usage if true.
	showHelp = flag.Bool("help", false, "")

	// Seconds between pings.
	spacing = flag.Int("spacing", 5, "")
)

const helpMessage = `

pingcheck is a command-line utility to ping ports and write log of responses.

Usage: pingcheck [options] <url1>, <url2> ...

      -spacing     =number  Seconds between pings.
  -h, -help       (flag)    Show help message
`

func main() {
	flag.BoolVar(showHelp, "h", false, "Show help message")
	flag.Usage = func() {
		fmt.Printf(helpMessage)
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(0)
	}

	log.Printf("Starting ping on %d targets...\n", flag.NArg())

	stopSig := make(chan os.Signal)
	signal.Notify(stopSig, os.Interrupt, os.Kill, syscall.SIGTERM)

	check := time.Tick(time.Duration(*spacing) * time.Second)
	for {
		select {
		case <-check:
			doPings()
		case <-stopSig:
			os.Exit(0)
		}
	}

	log.Printf("Halting pings.\n")
}

func doPings() {
	for _, url := range flag.Args() {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("ERROR on %q --> %v\n", url, err)
		} else {
			log.Printf("%q: %s\n", url, resp.Status)
		}
	}
}
