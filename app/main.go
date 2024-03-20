package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"syscall"
	"time"
)

var (
	logger  *log.Logger
	logFile *os.File
)

func init() {
	logFile, err := os.Create("logs/log.csv")
	if err != nil {
		log.Fatal("Failed To Create / Open logs/log.csv")
	}
	logger = log.New(logFile, "", 0)
}

func main() {
	//m := autocert.Manager{
	//	Prompt:     autocert.AcceptTOS,
	//	HostPolicy: autocert.HostWhitelist("npex42.dev"),
	//	Cache:      autocert.DirCache("certs"),
	//	Email:      "gvenn@npex42.dev",
	//}

	log.Default().SetFlags(0)

	http.Handle("/", RequestLogger(logger, http.FileServer(http.Dir("./static"))))
	http.Handle("/hooks/gh_push", RequestLogger(logger, WebhookPushHandler()))
	http.Handle("/hooks/pull", RequestLogger(logger, http.HandlerFunc(Restart)))

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT)

	go SignalLoop(sigChan)

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
			fmt.Printf("Starting Server @ %v\n", s.Addr)
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

func SignalLoop(sigChan chan os.Signal) {
	log.Println("[SignalHandler] - Started")
	for {
		SignalHandler(<-sigChan)
	}
}

func SignalHandler(sig os.Signal) {
	switch sig {
	case syscall.SIGINT:
		{
			log.Println("[SIGINT] Shutting Down...")
			os.Exit(0)
		}
	}
}

func backupLogs() {
	now := time.Now()
	copyPath := fmt.Sprintf("logs/log-%d%d%d.csv", now.Day(), now.Month(), now.Year())
	copyLog, err := os.Create(copyPath)
	if err == nil {
		io.Copy(copyLog, logFile)
	} else {
		log.Fatalln("Failed To Copy Log file...")
	}
}
