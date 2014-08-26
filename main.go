package main

import (
	"fmt"
	"flag"
)

func main() {
	fmt.Println("Launch succeeded!")

	confFile := flag.String("conf", "config.yml", "specify config file")
	flag.Parse()

	cfg := NewConfig(*confFile)

	Setup(cfg)
	Run()
}


