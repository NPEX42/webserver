package main

import (
	"fmt"
	"github.com/NPEX42/webserver/pkg/tools"
	"os"
)

type Frontmatter struct {
	Author string   `yaml:"Author"`
	Title  string   `yaml:"Title"`
	Tags   []string `yaml:"Tags"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <root Dir>\n", os.Args[0])
		os.Exit(1)
	}
	dir := os.Args[1]

	tools.PopulateBlogDB(dir)

}
