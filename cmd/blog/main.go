package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/NPEX42/webserver/app/blog"
	"github.com/adrg/frontmatter"
	"github.com/bmatcuk/doublestar/v4"
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
	matches, err := doublestar.Glob(os.DirFS(dir), "**/*.md")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d Posts\n", len(matches))

	os.Remove("./posts.db")
	pdb, err := blog.ConnectToPostDB("./posts.db")
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	pdb.Init()

	for _, match := range matches {
		file, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, match))
		if err != nil {
			log.Fatal(err)
		}
		var meta blog.PostMetadata
		rest, err := frontmatter.MustParse(strings.NewReader(string(file)), &meta)
		if err != nil {
			log.Printf("ERR IN POST %s: %s - SKIPPING\n", match, err)
			continue
		}

		fmt.Printf("%s - %s - %s\n", match, meta.Title, meta.Author)

		post := blog.Post{
			Metadata: meta,
			Content:  string(rest),
		}

		res, err := pdb.CreateNewPost(post)

		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}
		id, _ := (*res).LastInsertId()
		log.Printf("Created! ID: %d\n", id)
	}

}
