package main

import (
	"fmt"
	"flag"
	"log"
	"errors"
	"strings"
	"strconv"
	"github.com/acidlemon/mirage/mirage"
)

func main() {
	fmt.Println("Launch succeeded!")

	cfg := mirage.NewConfig()

	err := ParseArgs(cfg)
	if err != nil {
		log.Fatal(err)
	}

	mirage.Setup(cfg)
	mirage.Run()
}

func ParseArgs(cfg *mirage.Config) error {
	foreignAddress := flag.String("foreign-address", "127.0.0.1",
		"Listening foreign address")
	webApiHost := flag.String("webapi-host", "localhost",
		"Host name for webapi")
	domainSuffix := flag.String("domain-suffix", "",
		"domain suffix for reverse proxy. ex) \".example.com\"")
	endpoint := flag.String("endpoint", "unix:///var/run/docker.sock",
		"docker endpoint")
	listen := flag.String("listen", "8080", 
		"listen ports with mapping. you can specify multiple port using comma\nex) \"8080:5000,443\" ... listen 8080 and 443 port, and 8080 is proxy to 5000 port of backend")
	flag.Parse()

	if *domainSuffix == "" {
		return errors.New("-domain-suffix is required.")
	}

	cfg.ForeignAddress = *foreignAddress
	cfg.WebApiHost = *webApiHost
	cfg.ReverseProxyHostSuffix = *domainSuffix
	cfg.DockerEndpoint = *endpoint
	cfg.ListenPorts = ParseListenPort(*listen)

	return nil
}

func ParseListenPort(listenLine string) map[int]int {
	result := map[int]int{}

	portSets := strings.Split(listenLine, ",")
	for _, v := range portSets {
		ports := strings.Split(v, ":")
		var listen, target int
		listen, _ = strconv.Atoi(ports[0]) // ignore error because 0 is invalid
		if len(ports) > 1 {
			target, _ = strconv.Atoi(ports[1])
		} else {
			target, _ = strconv.Atoi(ports[0])
		}

		if isValidPort(listen) && isValidPort(target) {
			result[listen] = target
		} else {
			fmt.Println("invalid listen port: ", listenLine)
		}
	}

	return result
}

func isValidPort(port int) bool {
	return port > 0 && port < 65536
}

