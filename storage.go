package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

var ErrNotFound = errors.New("Not Found")

type MirageStorage struct {
	storage *leveldb.DB
}

func NewMirageStorage(cfg *Config) *MirageStorage {
	fileStorage, err := storage.OpenFile(cfg.Storage.DataDir)
	if err != nil {
		fmt.Println("cannot open leveldb fileStorage")
		log.Fatal(err)
	}

	storage, err := leveldb.Open(fileStorage, &opt.Options{})
	if err != nil {
		fmt.Println("cannot open leveldb")
		log.Fatal(err)
	}

	ms := &MirageStorage{storage: storage}

	return ms
}

func (ms *MirageStorage) Close() {
	ms.storage.Close()
}

func (ms *MirageStorage) Get(key string) ([]byte, error) {
	data, err := ms.storage.Get([]byte(key), nil)
	if err == errors.ErrNotFound {
		err = ErrNotFound
	}

	return data, err
}

func (ms *MirageStorage) Set(key string, value []byte) error {
	err := ms.storage.Put([]byte(key), value, &opt.WriteOptions{Sync: true})

	return err
}

func (ms *MirageStorage) AddToSubdomainMap(subdomain string) error {
	subdomainMap, err := ms.getSubdomainMap()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get subdomain-map: %s", err.Error()))
	}

	beforeLen := len(subdomainMap)

	subdomainMap[subdomain] = 1 // meanless value

	if beforeLen == len(subdomainMap) {
		// need not to update
		fmt.Println("subdomainMap length is not changed!")
		return nil
	}

	return ms.updateSubdomainMap(subdomainMap)
}

func (ms *MirageStorage) RemoveFromSubdomainMap(subdomain string) error {
	subdomainMap, err := ms.getSubdomainMap()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get subdomain-map: %s", err.Error()))
	}

	beforeLen := len(subdomainMap)

	delete(subdomainMap, subdomain)

	if beforeLen == len(subdomainMap) {
		return nil
	}

	return ms.updateSubdomainMap(subdomainMap)
}

func (ms *MirageStorage) getSubdomainMap() (map[string]int, error) {
	subdomainData, err := ms.Get("subdomain-map")
	if err != nil {
		if err != ErrNotFound {
			return nil, err
		}
		subdomainData = []byte(`{}`)
	}

	// Q. Why map?  A. It's easy to manage subdomains as unique
	var subdomainMap map[string]int
	err = json.Unmarshal(subdomainData, &subdomainMap)
	if err != nil {
		return nil, err
	}

	return subdomainMap, nil
}

func (ms *MirageStorage) GetSubdomainList() ([]string, error) {
	subdomainMap, err := ms.getSubdomainMap()
	if err != nil {
		return []string{}, err
	}

	result := []string{}
	for k, _ := range subdomainMap {
		result = append(result, k)
	}

	return result, nil
}

func (ms *MirageStorage) updateSubdomainMap(subdomainMap map[string]int) error {
	//dump.Dump(subdomainMap)
	subdomainData, err := json.Marshal(subdomainMap)

	err = ms.Set("subdomain-map", subdomainData)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to update subdomain-map: %s", err.Error()))
	}

	return nil
}
