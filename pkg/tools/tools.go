package tools

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/NPEX42/webserver/app/blog"
	"github.com/adrg/frontmatter"
	"github.com/bmatcuk/doublestar/v4"
)

func PopulateBlogDB(dir string) {
	matches, err := doublestar.Glob(os.DirFS(dir), "**/*.md")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d Posts\n", len(matches))

	pdb, err := blog.ConnectToPostDB("./posts.db")
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	pdb.Init()

	defer pdb.Close()

	for _, match := range matches {

		filepath := fmt.Sprintf("%s/%s", dir, match)

		if strings.Contains(filepath, "excalidraw") {
			continue
		}

		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}

		file, err := os.ReadFile(filepath)
		if err != nil {
			log.Fatal(err)
		}
		var meta blog.PostMetadata

		hash := sha1.New()

		uploadTime := time.Now()
		var b []byte

		b = binary.BigEndian.AppendUint16(b, uint16(uploadTime.UnixMilli()))

		if _, err := io.Copy(hash, bytes.NewReader(b)); err != nil {
			log.Fatal(err)
		}

		_, err = frontmatter.MustParse(strings.NewReader(string(file)), &meta)
		if err != nil {
			log.Printf("ERR IN POST %s: %s - SKIPPING\n", match, err)
			continue
		}
		meta.ID = base64.RawURLEncoding.EncodeToString(hash.Sum(nil)[0:4])

		post := blog.Post{
			Metadata: meta,
			Filepath: filepath,
		}

		samePosts, err := pdb.GetPostByFilepath(post.Filepath)

		log.Println(samePosts)
		if len(samePosts) == 0 {
			_, err = pdb.CreateNewPost(post)
			if err != nil {
				log.Fatalf("Error: %s\n", err)
			}
			fmt.Printf("Created: %s - %s - %s - %s\n", match, meta.Title, meta.Author, meta.ID)
		} else {
			log.Println("Skipping Post")
			continue
		}

		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}
	}

}

func DeletePost(id string) {
	pdb, err := blog.ConnectToPostDB("./posts.db")
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	defer pdb.Close()

	err = pdb.DeletePostByID(id)
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
}
