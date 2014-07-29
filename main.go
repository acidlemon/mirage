package main

import (
	"fmt"
	"github.com/acidlemon/mirage/mirage"
)

func main() {
	fmt.Println("Launch succeeded!")

	cfg := mirage.NewConfig()
	mirage.Setup(cfg)
	mirage.Run()
}

