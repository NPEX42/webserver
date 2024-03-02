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

func main() {
	//m := autocert.Manager{
	//	Prompt:     autocert.AcceptTOS,
	//	HostPolicy: autocert.HostWhitelist("npex42.dev"),
	//	Cache:      autocert.DirCache("certs"),
	//	Email:      "gvenn@npex42.dev",
	//}

	log.Default().SetFlags(0)

	http.Handle("/", RequestLogger(log.Default(), http.FileServer(http.Dir("./static"))))
	http.Handle("/hooks/gh_push", RequestLogger(log.Default(), WebhookPushHandler()))

	conf, err := LoadServerConfig("config.json")
	if err != nil {
		log.Fatal(err)
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
			log.Printf("Starting Server @ %v\n", s.Addr)
			log.Fatal(s.ListenAndServe())
		}
		return
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
	CertDir   string `json:"certDir"`
	AllowHTTP bool   `json:"allowHTTP"`
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

type WrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *WrappedResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &WrappedResponseWriter{w, http.StatusOK}
}

func (wrw *WrappedResponseWriter) WriteHeader(code int) {
	wrw.statusCode = code
	wrw.ResponseWriter.WriteHeader(code)
}
