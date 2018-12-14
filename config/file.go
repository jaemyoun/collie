package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type configFileType struct {
	Buckets []string               `json:"buckets"`
	Filters []ConfigFileFilterType `json:"filters"`
}

type ConfigFileFilterType struct {
	Expr string `json:"expr"`
	Desc string `json:"desc"`
}

var Predefined configFileType

func readConfigFile() {
	body, err := ioutil.ReadFile(os.Getenv("HOME") + "/.collie")
	if err != nil {
		log.Println("there is no config file")
		return
	}

	err = json.Unmarshal(body, &Predefined)
	if err != nil {
		log.Fatal("cannot unmarshal config file")
	}

	for idx, e := range Predefined.Filters {
		Predefined.Filters[idx].Expr = strings.Replace(e.Expr, `\`, `\\`, -1)
	}
}

func (c configFileType) GetBuckets() []string {
	return c.Buckets
}

func (c configFileType) GetFilters() []ConfigFileFilterType {
	return c.Filters
}
