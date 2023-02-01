package client

import (
	"fmt"
	"os"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
	uuid "github.com/nu7hatch/gouuid"
	"gopkg.in/yaml.v2"
)

// Config configures an Proxy
type Config struct {
	id           string
	subdomain    string
	port         int
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
	config.id = id.String()

	config.Targets = []string{"wss://t.amal.sh"}
	config.PoolIdleSize = 1
	config.PoolMaxSize = 100

	return
}

// LoadConfiguration loads configuration from a YAML file
func LoadConfiguration(path, subdomain string, port int) (config *Config, err error) {
	config = NewConfig()

	bytes, err := os.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return
	}
	for i, v := range config.Targets {
		if !strings.HasSuffix(v, "/register") {
			config.Targets[i] = fmt.Sprintf("%s/register", v)
		}
	}

	if subdomain == "" {
		subdomain, err = gonanoid.Generate("abcdefghijklmnopqrstuvwxyz", 6)
		if err != nil {
			panic(err)
		}
	}
	config.subdomain = subdomain
	config.port = port

	return
}
