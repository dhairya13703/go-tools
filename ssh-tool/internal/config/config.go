package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

//go:embed servers.json
var embeddedConfig []byte

type Server struct {
	Name        string
	Hostname    string `json:"hostname"`
	User        string `json:"user"`
	PemFile     string `json:"pem_file"`
	Description string `json:"description"`
}

type Config struct {
	Servers map[string]Server `json:"servers"`
}

func LoadConfig(file string) (*Config, error) {
	var data []byte
	var err error

	if file != "" {
		// Try to load from external file if provided
		data, err = os.ReadFile(file)
		if err != nil {
			// If external file fails, fall back to embedded config
			data = embeddedConfig
		}
	} else {
		// Use embedded config by default
		data = embeddedConfig
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config: %v", err)
	}

	// Add server names to the struct
	for name, server := range config.Servers {
		server.Name = name
		config.Servers[name] = server
	}

	return &config, nil
}

func (c *Config) GetServersList() []Server {
	servers := make([]Server, 0, len(c.Servers))
	for name, server := range c.Servers {
		server.Name = name
		servers = append(servers, server)
	}

	// Sort servers by name for consistent ordering
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Name < servers[j].Name
	})

	return servers
}
