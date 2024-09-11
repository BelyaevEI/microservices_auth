package config

import (
	"errors"
	"net"
	"os"
)

const (
	httpHostEnvName = "HOST_PROMETHEUS"
	httpPortEnvName = "PORT_PROMETHEUS"
)

// PrometheusConfig config for prometheus server
type PrometheusConfig interface {
	Address() string
}

type prometheusConfig struct {
	host string
	port string
}

// NewPrometheusConfig initializes a prometheus configuration.
func NewPrometheusConfig() (PrometheusConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("http host not found")
	}

	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("http port not found")
	}

	return &prometheusConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *prometheusConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
