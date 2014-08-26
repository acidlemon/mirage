package main

import (
	"log"
	"io/ioutil"
	"strings"
	"fmt"
	"strconv"
	"gopkg.in/yaml.v1"
)

type Config struct {
	Host Host                  `yaml:"host"`
	Listen Listen              `yaml:"listen"`
	Docker DockerCfg           `yaml:"docker"`
}

type Host struct {
	WebApi string              `yaml:"webapi"`
	ReverseProxySuffix string  `yaml:"reverse_proxy_suffix"`
}

type Listen struct {
	ForeignAddress string      `yaml:"foreign_address"`
	http []string              `yaml:"http"`
	https []string             `yaml:"https"`
	HTTP map[int]int
	HTTPS map[int]int
}

type DockerCfg struct {
	Endpoint string            `yaml:"endpoint"`
	DefaultImage string        `yaml:"default_image"`
}

func NewConfig(path string) *Config {
	// default config
	cfg := &Config{
		Host: Host{
			WebApi: "localhost",
			ReverseProxySuffix: ".dev.example.net",
		},
		Listen: Listen{
			ForeignAddress: "127.0.0.1",
			HTTP: map[int]int{ 8080: 5000 },
			HTTPS: map[int]int{},
		},
		Docker: DockerCfg{
			Endpoint: "unix:///var/run/docker.sock",
			DefaultImage: "",
		},
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("cannot read %v: %v", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		log.Fatalf("%v", err)
	}

	cfg.Listen.Parse()

	return cfg
}

func (l Listen) Parse(){
	l.HTTP = parseListenPort(l.http)
	l.HTTPS = parseListenPort(l.https)
}

func parseListenPort(portSets []string) map[int]int {
	result := map[int]int{}

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
			fmt.Println("invalid listen port: ", v)
		}
	}

	return result
}

func isValidPort(port int) bool {
	return port > 0 && port < 65536
}


