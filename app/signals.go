package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

var signal_logger log.Logger

func StartSignalHandler() {
	signal_logger = *log.New(os.Stdout, "[Signal Handler] - ", log.Lmsgprefix)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	SignalLoop(sigChan)
}

func SignalLoop(sigChan chan os.Signal) {
	signal_logger.Println("Started")
	for {
		SignalHandler(<-sigChan)
	}
}

func SignalHandler(sig os.Signal) {
	switch sig {
	case syscall.SIGINT:
		{
			signal_logger.Println("[SIGINT] Shutting Down...")
			os.Exit(0)
		}
	}
}
