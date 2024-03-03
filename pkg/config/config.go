package config

import (
	"os"

	"gopkg.in/yaml.v1"
)

type Config struct {
	BufSizePerProc int       `yaml:"bufSizePerProc"`
	NumOfProcs     int       `yaml:"numOfProcs"`
	TCPConns       []TCPConn `yaml:"tcpConns"`
	UDPConns       []UDPConn `yaml:"udpConns"`
}

type TCPConn struct {
	SrcAddr string `yaml:"srcAddr"`
	DstAddr string `yaml:"dstAddr"`
}

type UDPConn struct {
	SrcAddr string `yaml:"srcAddr"`
	DstAddr string `yaml:"dstAddr"`
}

func NewConfig() *Config {
	return &Config{}
}

func ParseConfigFile(configPath string) (*Config, error) {
	config := NewConfig()
	content, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		return config, err
	}

	return config, nil
}
