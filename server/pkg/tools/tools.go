package tools

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/NPEX42/webserver/blog/pkg/blog"
	"github.com/bmatcuk/doublestar/v4"
)

func PopulateBlogDB(dir string) {
	matches, err := doublestar.Glob(os.DirFS(dir), "**/*.md")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d Posts\n", len(matches))

	pdb, err := blog.Connect("./posts.json")
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	defer pdb.Save()

	for _, match := range matches {

		filepath := fmt.Sprintf("%s/%s", dir, match)

		if strings.Contains(filepath, "excalidraw") {
			continue
		}

		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}

		post, err := blog.LoadFromMarkdownFile(filepath)

		post.ID = blog.GetPostID(time.Now())

		samePosts, err := pdb.GetPostByFilePath(post.Filepath)

		log.Println(samePosts)
		if len(samePosts) == 0 {
			err = pdb.CreateNewPost(post)
			if err != nil {
				log.Fatalf("Error: %s\n", err)
			}
			fmt.Printf("Created: %s - %s - %s - %s\n", match, post.Title, post.Author, post.ID)
		} else {
			log.Println("Skipping Post")
			continue
		}

		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}
	}

	pdb.Save()

}

func DeletePost(id string) {
	pdb, err := blog.Connect("./posts.json")
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	pdb.Save()

	err = pdb.DeletePostByID(id)
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
}
