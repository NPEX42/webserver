package main

import (
	"fmt"
	"os"

	"github.com/NPEX42/webserver/pkg/tools"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <root Dir>\n", os.Args[0])
		os.Exit(1)
	}
	id := os.Args[1]
	tools.DeletePost(id)
}
