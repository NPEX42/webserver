package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NPEX42/webserver/blog/pkg/blog"
)

var signal_logger log.Logger

func StartSignalHandler(pdb *blog.Blog) {
	signal_logger = *log.New(os.Stdout, "[Signal Handler] - ", log.Lmsgprefix)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	SignalLoop(sigChan, pdb)
}

func SignalLoop(sigChan chan os.Signal, pdb *blog.Blog) {
	signal_logger.Println("Started")
	for {
		SignalHandler(<-sigChan, pdb)
	}
}

func SignalHandler(sig os.Signal, pdb *blog.Blog) {
	switch sig {
	case syscall.SIGINT:
		{
			signal_logger.Println("[SIGINT] Shutting Down...")
			pdb.Save()
			os.Exit(0)
		}
	}
}
