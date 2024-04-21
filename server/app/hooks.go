package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"syscall"
)

func Restart(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Pulling...")
	Pull()
	exec.Command("go", "build", "-o", "wserver", "./app").Run()
	SignalHandler(syscall.SIGHUP)
	syscall.Exec("./wserver", []string{}, []string{})
}
