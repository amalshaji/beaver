package client

import (
	"fmt"
	"os"
	"strings"

	"github.com/amalshaji/beaver/internal/utils"
	gonanoid "github.com/matoous/go-nanoid/v2"
	uuid "github.com/nu7hatch/gouuid"

	"gopkg.in/yaml.v3"
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
func LoadConfiguration(
	configFile string,
	subdomain string,
	port int,
	showWsReadErrors bool,
) (Config, error) {
	var config Config

	bytes, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return Config{}, err
	}

	config.setDefaults()

	if !strings.HasSuffix(config.Target, "/register") {
		config.Target = config.Target + "/register"
	}

	if subdomain == "" {
		subdomain, err = gonanoid.Generate("abcdefghijklmnopqrstuvwxyz", 6)
		if err != nil {
			panic(err)
		}
	} else {
		err = utils.ValidateSubdomain(subdomain)
		if err != nil {
			return Config{}, fmt.Errorf("invalid subdomain: '%s'; %s", subdomain, err.Error())
		}
	}

	config.subdomain = subdomain
	config.port = port
	config.showWsReadErrors = showWsReadErrors

	return config, nil
}
