package config

import (
	sharedcfg "github.com/Muxly-Corp/muxly-shared/config"
)

type Config struct {
	HTTP         sharedcfg.HTTPConfig
	Log          sharedcfg.LogConfig
	Queue        sharedcfg.QueueConfig
	MuxlyBackend MuxlyBackend
	Subscriber   Subscriber
}

type MuxlyBackend struct {
	URL            string `env:"MUXLY_BACKEND_URL"              envDefault:"http://localhost:3001"`
	InternalAPIKey string `env:"MUXLY_BACKEND_INTERNAL_API_KEY"`
}

type Subscriber struct {
	BatchIntervalMs int `env:"SUBSCRIBER_BATCH_INTERVAL_MS" envDefault:"100"`
	BatchMaxSize    int `env:"SUBSCRIBER_BATCH_MAX_SIZE"    envDefault:"50"`
}

// Mandatory AppConfig interface methods.
func (c *Config) GetHTTPConfig() sharedcfg.HTTPConfig { return c.HTTP }
func (c *Config) GetLogConfig() sharedcfg.LogConfig   { return c.Log }

// Optional infra interface methods.
func (c *Config) GetQueueConfig() sharedcfg.QueueConfig { return c.Queue }
