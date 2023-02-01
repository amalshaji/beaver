package client

import (
	"os"

	uuid "github.com/nu7hatch/gouuid"
	"gopkg.in/yaml.v2"
)

// Config configures an Proxy
type Config struct {
	ID           string
	Targets      []string
	PoolIdleSize int
	PoolMaxSize  int
	SecretKey    string
}

// NewConfig creates a new ProxyConfig
func NewConfig() (config *Config) {
	config = new(Config)

	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	config.ID = id.String()

	config.Targets = []string{"wss://t.amal.sh"}
	config.PoolIdleSize = 1
	config.PoolMaxSize = 100

	return
}

// LoadConfiguration loads configuration from a YAML file
func LoadConfiguration(path string) (config *Config, err error) {
	config = NewConfig()

	bytes, err := os.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return
	}

	return
}
