package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

const (
	ENV_ENDPOINT_HOST = "ENDPOINT_HOST"
	ENV_ENDPOINT_PORT = "ENDPOINT_PORT"
)

type AppCfg struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type LogCfg struct {
	Location string `json:"location"`
	Service  string `json:"service"`
	File     string `json:"file"`
	Level    string `json:"level"`
}

type AWSCfg struct {
	Region    string  `json:"region"`
	Host      string  `json:"host"`
	Port      int     `json:"port"`
	AccountId string  `json:"account_id"`
	Dev       bool    `json:"dev"`
	SQS       *SQSCfg `json:"sqs"`
}

type SQSCfg struct {
	Name              string  `json:"name"`
	IsFifo            bool    `json:"is_fifo"`
	VisibilityTimeout int64   `json:"visibility_timeout"`
	ContextTimeout    int64   `json:"context_timeout"`
	Retention         int64   `json:"retention"`
	PollInterval      int     `json:"poll_interval"`
	Log               *LogCfg `json:"log"`
}

type Config struct {
	App *AppCfg `json:"app"`
	Log *LogCfg `json:"log"`
	AWS *AWSCfg `json:"aws"`
}

type ConfigInput struct {
	Name string
	Type string
	Path string
}

func GetConfig(cin *ConfigInput) *Config {
	viper.Reset()

	viper.SetConfigName(cin.Name)
	viper.SetConfigType(cin.Type)
	viper.AddConfigPath(cin.Path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("could not read config file at %v: %v", cin, err)
		return nil
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatalf("could not parse config file at %v: %v", cin, err)
		return nil
	}

	// populate endpoint from env var
	endpointHost := os.Getenv(ENV_ENDPOINT_HOST)
	if endpointHost != "" {
		cfg.AWS.Host = endpointHost
	}
	endpointPort := os.Getenv(ENV_ENDPOINT_PORT)
	if endpointPort != "" {
		cfg.AWS.Port, _ = strconv.Atoi(endpointPort)
	}

	return cfg
}
