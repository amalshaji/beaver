package client

import (
	"fmt"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
	uuid "github.com/nu7hatch/gouuid"
)

type TunnelConfig struct {
	Subdomain string
	Port      int
}

type Config struct {
	Targets      []string
	PoolIdleSize int
	PoolMaxSize  int
	SecretKey    string
}

// Proxy configures an Proxy
type Proxy struct {
	id               string
	subdomain        string
	port             int
	showWsReadErrors bool

	Config  Config
	Tunnels []TunnelConfig
}

func (proxy *Proxy) setDefaults() {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	proxy.id = id.String()

	if len(proxy.Config.Targets) == 0 {
		proxy.Config.Targets = []string{"wss://x.amal.sh"}
	}
	if proxy.Config.PoolIdleSize == 0 {
		proxy.Config.PoolIdleSize = 1
	}
	if proxy.Config.PoolMaxSize == 0 {
		proxy.Config.PoolMaxSize = 100
	}

}

// LoadConfiguration loads configuration from a YAML file
func LoadConfiguration(config Proxy, subdomain string, port int, showWsReadErrors bool) (Proxy, error) {
	var err error

	config.setDefaults()

	for i, v := range config.Config.Targets {
		if !strings.HasSuffix(v, "/register") {
			config.Config.Targets[i] = fmt.Sprintf("%s/register", v)
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
	config.showWsReadErrors = showWsReadErrors

	config.Tunnels = make([]TunnelConfig, 0)

	return config, nil
}
