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
	exec.Command("go", "build", "-o", "server", "./app").Run()
	syscall.Exec("./server", []string{}, []string{})
}
