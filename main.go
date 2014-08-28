package main

import (
	"fmt"
	"flag"
)

var (
	version string
	buildDate string
)

func main() {
	confFile := flag.String("conf", "config.yml", "specify config file")
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.Parse()

	if showVersion {
		fmt.Printf("mirage %v (%v)\n", version, buildDate)
		return
	}

	fmt.Println("Launch succeeded!")

	cfg := NewConfig(*confFile)

	Setup(cfg)
	Run()
}


