package config

import (
	"github.com/spf13/viper"
	"log"
)

var (
	AppConfig     *Config
	AppConfigPath string
)

type Server struct {
	Listen string
}

type Proxy struct {
	Location  string
	ProxyPass string
}

type Config struct {
	Server Server
	Proxy  []Proxy
}

func Init(path string) {
	cfg, err := LoadConfig(path)
	if err != nil {
		log.Fatalf("Load Config error: %v", err)
	}
	config, err := ParseConfig(cfg)
	if err != nil {
		log.Fatalf("Parse Config error: %v", err)
	}

	AppConfig = config
}

func LoadConfig(path string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (c *Config, err error) {
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}
	return c, nil
}
