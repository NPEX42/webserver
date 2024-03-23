package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"time"
)

var (
	logger  *log.Logger
	logFile *os.File
)

func init() {
	now := time.Now()
	logFile, err := os.Create(fmt.Sprintf("logs/log-%d%d%d.csv", now.Day(), now.Month(), now.Year()))
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

	go StartSignalHandler()

	http.Handle("/hooks/gh_push", RequestLogger(logger, WebhookPushHandler()))
	http.Handle("/hooks/pull", RequestLogger(logger, http.HandlerFunc(Restart)))
	http.Handle("/api/v1/projects", RequestLogger(logger, http.HandlerFunc(GetProjects)))

	conf, err := LoadServerConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", RequestLogger(logger, http.FileServer(http.Dir(conf.StaticDir))))

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

func backupLogs() {
	now := time.Now()
	copyPath := fmt.Sprintf("logs/log-%d%d%d.csv", now.Day(), now.Month(), now.Year())
	copyLog, err := os.Create(copyPath)
	if err == nil {
		io.Copy(copyLog, logFile)
		copyLog.Close()
	} else {
		log.Fatalln("Failed To Copy Log file...")
	}
}
