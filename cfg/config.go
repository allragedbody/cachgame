package cfg

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Server       *Server
	Organization *CfgOrganization
	Country      *CfgCountry
}

type Server struct {
	Port  int
	Cores int
}

type CfgOrganization struct {
	Name  string
	Retio float64
}

type CfgCountry struct {
	Name  string
	Retio float64
}

func ParseConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
