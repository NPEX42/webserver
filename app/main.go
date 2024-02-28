package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	//m := autocert.Manager{
	//	Prompt:     autocert.AcceptTOS,
	//	HostPolicy: autocert.HostWhitelist("npex42.dev"),
	//	Cache:      autocert.DirCache("certs"),
	//	Email:      "gvenn@npex42.dev",
	//}

	http.Handle("/", RequestLogger(log.Default(), http.FileServer(http.Dir("./static"))))

	conf, err := LoadServerConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	cert, err := conf.LoadCertificate()
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	s.ListenAndServeTLS("", "")
}

type ServerConfig struct {
	CertDir string `json:"certDir"`
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
	certFile := fmt.Sprintf("%s/cert.pem", conf.CertDir)
	keyFile := fmt.Sprintf("%s/privkey.pem", conf.CertDir)

	return tls.LoadX509KeyPair(certFile, keyFile)
}
