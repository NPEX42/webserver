package main

import (
	"context"
	"log"
	"os"

	"github.com/NPEX42/webserver/components/common"
)

func main() {
	file, err := os.Create("static/home.html")
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	err = common.Header("Home").Render(context.Background(), file)

	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

}
