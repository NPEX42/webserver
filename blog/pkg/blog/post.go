package blog

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Post struct {
	ID       string `json:"ID"`
	Filepath string `json:"filepath"`
	Author   string `yaml:"Author" json:"Author"`
	Title    string `yaml:"Title" json:"Title"`
}

func GetPostID(now time.Time) string {
	utc := now.UTC()
	_ = utc
	var b []byte

	b = binary.BigEndian.AppendUint16(b, uint16(utc.UnixMilli()))

	hash := sha1.New()

	if _, err := io.Copy(hash, bytes.NewReader(b)); err != nil {
		log.Fatal(err)
	}

	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil)[0:4])
}

func (post *Post) RenderHTML() string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)

	contents, err := os.ReadFile(post.Filepath)

	if err != nil {
		return "<h3>Failed To Render Post</h3>"
	}

	var _meta Post

	rest, err := frontmatter.Parse(strings.NewReader(string(contents)), &_meta)

	doc := p.Parse(rest)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.LazyLoadImages

	opts := html.RendererOptions{
		Flags: htmlFlags,
	}
	renderer := html.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

func LoadFromMarkdown(md string) (Post, error) {
	var post Post
	_, err := frontmatter.MustParse(strings.NewReader(md), &post)
	if err != nil {
		return post, err
	}

	return post, nil
}

func LoadFromMarkdownFile(filepath string) (Post, error) {
	md, err := os.ReadFile(filepath)
	var post Post
	if err != nil {
		return post, err
	}

	post, err = LoadFromMarkdown(string(md))
	post.Filepath = filepath
	return post, nil
}
