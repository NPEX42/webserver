package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"

	"github.com/labstack/echo/v4"
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

func (conf *ServerConfig) LoadCertificatePaths() (string, string) {
	certFile := fmt.Sprintf("%s/fullchain.pem", conf.CertDir)
	keyFile := fmt.Sprintf("%s/privkey.pem", conf.CertDir)

	return certFile, keyFile
}

func StartWebserver(path string, e *echo.Echo) error {

	s, err := CreateWebserver(path)

	if err != nil {
		return err
	}

	e.StartServer(s)

	return nil

}

func CreateWebserver(path string) (*http.Server, error) {

	conf, err := LoadServerConfig(path)
	if err != nil {
		return nil, err
	}

	cert, err := conf.LoadCertificate()
	if err != nil {
		log.Printf("Failed To Locate TLS Certificates, Falling back to HTTP.\n Error: %v\n", err)
		currentUser, err := user.Current()
		if err != nil {
			log.Fatalf("Failed To Get Current User. %v\n", err)
		}

		var s *http.Server

		if currentUser.Username == "root" {
			s = &http.Server{
				Addr: ":http",
			}
		} else {
			s = &http.Server{
				Addr: ":8080",
			}
		}
		if conf.AllowHTTP {
			fmt.Printf("Starting Server @ %v\n", s.Addr)
			return s, nil
		}
		return nil, fmt.Errorf("Failed To Create Server\n")
	}

	s := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	return s, nil
}
