package client

import (
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
	uuid "github.com/nu7hatch/gouuid"
)

type TunnelConfig struct {
	Name      string
	Subdomain string
	Port      int
}

type ProxyTunnels struct {
	Tunnels []TunnelConfig
}

type Config struct {
	id               string
	subdomain        string
	port             int
	showWsReadErrors bool

	Target       string
	PoolIdleSize int
	PoolMaxSize  int
	SecretKey    string
}

// ProxyConfig configures an ProxyConfig

func (config *Config) setDefaults() {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	config.id = id.String()

	if config.Target == "" {
		config.Target = "wss://x.amal.sh"
	}
	if config.PoolIdleSize == 0 {
		config.PoolIdleSize = 1
	}
	if config.PoolMaxSize == 0 {
		config.PoolMaxSize = 100
	}

}

// LoadConfiguration loads configuration from a YAML file
func LoadConfiguration(config Config, subdomain string, port int, showWsReadErrors bool) (Config, error) {
	var err error

	config.setDefaults()

	if !strings.HasSuffix(config.Target, "/register") {
		config.Target = config.Target + "/register"
	}

	if subdomain == "" {
		subdomain, err = gonanoid.Generate("abcdefghijklmnopqrstuvwxyz", 6)
		if err != nil {
			panic(err)
		}
	}
	config.subdomain = subdomain
	config.port = port
	config.showWsReadErrors = showWsReadErrors

	return config, nil
}
