package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NPEX42/webserver/app/blog"
	"github.com/NPEX42/webserver/pkg/tools"
	"github.com/labstack/echo/v4"
)

var (
	logger  *log.Logger
	logFile *os.File
	pdb     *blog.PostDB
)

func init() {
	now := time.Now()
	logFile, err := os.Create(fmt.Sprintf("logs/log-%d%d%d.csv", now.Day(), now.Month(), now.Year()))
	if err != nil {
		log.Fatal("Failed To Create / Open logs/log.csv")
	}

	logger = log.New(logFile, "", 0)
	logger.Println("Date-Time,Method,RemoteIP,Path,Status,DurationMicro")

	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	tools.PopulateBlogDB("./posts")

	pdb, err := blog.ConnectToPostDB("./posts.db")
	if err != nil {
		log.Fatalf("Failed To Connect To Posts Database: %s\n", err)
	}

	log.Default().SetFlags(0)

	go StartSignalHandler()
	conf, err := LoadServerConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	CreateEcho(&conf, e, pdb)
	err = StartWebserver("config.json", e)
	if err != nil {
		log.Fatalln(err)
	}

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
