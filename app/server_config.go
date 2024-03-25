package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
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

func StartWebserver(path string, router *http.ServeMux) error {

	conf, err := LoadServerConfig(path)
	if err != nil {
		return err
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
		s.Handler = router
		if conf.AllowHTTP {
			fmt.Printf("Starting Server @ %v\n", s.Addr)
			log.Fatal(s.ListenAndServe())
		}
		return nil
	}

	s := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
		Handler: router,
	}

	return s.ListenAndServeTLS("", "")

}
