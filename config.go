package main

import (
	"fmt"
	"github.com/acidlemon/go-dumper"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
)

type Config struct {
	Host   Host      `yaml:"host"`
	Listen Listen    `yaml:"listen"`
	Docker DockerCfg `yaml:"docker"`
}

type Host struct {
	WebApi             string `yaml:"webapi"`
	ReverseProxySuffix string `yaml:"reverse_proxy_suffix"`
}

type Listen struct {
	ForeignAddress string    `yaml:"foreign_address"`
	HTTP           []PortMap `yaml:"http"`
	HTTPS          []PortMap `yaml:"https"`
}

type PortMap struct {
	ListenPort int `yaml:"listen"`
	TargetPort int `yaml:"target"`
}

type DockerCfg struct {
	Endpoint     string `yaml:"endpoint"`
	DefaultImage string `yaml:"default_image"`
}

func NewConfig(path string) *Config {
	// default config
	cfg := &Config{
		Host: Host{
			WebApi:             "localhost",
			ReverseProxySuffix: ".dev.example.net",
		},
		Listen: Listen{
			ForeignAddress: "127.0.0.1",
			HTTP:           []PortMap{},
			HTTPS:          []PortMap{},
		},
		Docker: DockerCfg{
			Endpoint:     "unix:///var/run/docker.sock",
			DefaultImage: "",
		},
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("cannot read %v: %v", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		log.Fatalf("powawa: %v", err)
	}

	fmt.Println("Config:")
	dump.Dump(cfg)

	return cfg
}
