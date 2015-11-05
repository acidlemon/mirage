package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/fsouza/go-dockerclient"
)

type Information struct {
	ID        string `json:"id"`
	ShortID   string `json:"short_id"`
	SubDomain string `json:"subdomain"`
	GitBranch string `json:"branch"`
	Image     string `json:"image"`
	IPAddress string `json:"ipaddress"`
}

type Docker struct {
	cfg     *Config
	Storage *MirageStorage
	Client  *docker.Client
}

func NewDocker(cfg *Config, ms *MirageStorage) *Docker {
	client, err := docker.NewClient(cfg.Docker.Endpoint)
	if err != nil {
		fmt.Println("cannot create docker client")
		return nil
	}
	d := &Docker{
		cfg:     cfg,
		Storage: ms,
		Client:  client,
	}

	return d
}

func (d *Docker) Launch(subdomain string, image string, option map[string]string) error {
	var dockerEnv []string = make([]string, 0)
	for _, v := range d.cfg.Parameter {
		if option[v.Name] == "" {
			continue
		}

		dockerEnv = append(dockerEnv, fmt.Sprintf("%s=%s", v.Env, option[v.Name]))
	}
	dockerEnv = append(dockerEnv, fmt.Sprintf("SUBDOMAIN=%s", subdomain))

	opt := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: image,
			Env:   dockerEnv,
		},
	}

	container, err := d.Client.CreateContainer(opt)
	if err != nil {
		fmt.Println("cannot create container")
		return err
	}

	err = d.Client.StartContainer(container.ID, nil)
	if err != nil {
		fmt.Println("cannot start container")
		return err
	}

	container, err = d.Client.InspectContainer(container.ID)

	ms := d.Storage

	// terminate old container
	oldContainerID := d.getContainerIDFromSubdomain(subdomain, ms)
	if oldContainerID != "" {
		err = d.Client.StopContainer(oldContainerID, 5)
		if err != nil {
			fmt.Printf(err.Error()) // TODO log warning
		}
	}

	info := Information{
		ID:        container.ID,
		ShortID:   container.ID[0:12],
		SubDomain: subdomain,
		GitBranch: option["branch"],
		Image:     image,
		IPAddress: container.NetworkSettings.IPAddress,
	}
	var infoData []byte
	infoData, err = json.Marshal(info)

	err = ms.Set(fmt.Sprintf("subdomain:%s", subdomain), infoData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ms.AddToSubdomainMap(subdomain)
	app.ReverseProxy.AddSubdomain(subdomain, container.NetworkSettings.IPAddress)

	return nil
}

func (d *Docker) getContainerIDFromSubdomain(subdomain string, ms *MirageStorage) string {
	data, err := ms.Get(fmt.Sprintf("subdomain:%s", subdomain))
	if err != nil {
		if err == ErrNotFound {
			return ""
		}
		fmt.Printf("cannot find subdomain:%s, err:%s", subdomain, err.Error())
		return ""
	}
	var info Information
	json.Unmarshal(data, &info)
	//dump.Dump(info)
	containerID := string(info.ID)

	return containerID
}

func (d *Docker) Terminate(subdomain string) error {
	ms := d.Storage

	containerID := d.getContainerIDFromSubdomain(subdomain, ms)

	err := d.Client.StopContainer(containerID, 5)
	if err != nil {
		return err
	}

	ms.RemoveFromSubdomainMap(subdomain)
	app.ReverseProxy.RemoveSubdomain(subdomain)

	return nil
}

// extends docker.APIContainers for sort pkg
type ContainerSlice []docker.APIContainers

func (c ContainerSlice) Len() int {
	return len(c)
}
func (c ContainerSlice) Less(i, j int) bool {
	return c[i].ID < c[j].ID
}
func (c ContainerSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (d *Docker) List() ([]Information, error) {
	ms := d.Storage
	subdomainList, err := ms.GetSubdomainList()
	if err != nil {
		return nil, err
	}

	containers, _ := d.Client.ListContainers(docker.ListContainersOptions{})
	sort.Sort(ContainerSlice(containers))

	result := []Information{}
	for _, subdomain := range subdomainList {
		infoData, err := ms.Get(fmt.Sprintf("subdomain:%s", subdomain))
		if err != nil {
			fmt.Printf("ms.Get failed err=%s\n", err.Error())
			continue
		}

		var info Information
		err = json.Unmarshal(infoData, &info)
		//dump.Dump(info)

		index := sort.Search(len(containers), func(i int) bool { return containers[i].ID >= info.ID })

		if index < len(containers) && containers[index].ID == info.ID {
			// found
			result = append(result, info)
		}
	}

	return result, nil
}
