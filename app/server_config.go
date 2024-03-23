package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
)

type ServerConfig struct {
	CertDir   string `json:"certDir"`
	AllowHTTP bool   `json:"allowHTTP"`
	StaticDir string `json:"staticDir"`
}

func LoadServerConfig(path string) (ServerConfig, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return ServerConfig{}, err
	}
	var config ServerConfig
	err = json.Unmarshal(source, &config)
	if err != nil {
		return ServerConfig{}, err
	}

	return config, nil
}

func (conf *ServerConfig) LoadCertificate() (tls.Certificate, error) {
	certFile := fmt.Sprintf("%s/fullchain.pem", conf.CertDir)
	keyFile := fmt.Sprintf("%s/privkey.pem", conf.CertDir)

	return tls.LoadX509KeyPair(certFile, keyFile)
}
