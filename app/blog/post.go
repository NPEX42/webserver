package blog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type PostMetadata struct {
	Title  string `yaml:"Title"`
	Author string `yaml:"Author"`
	ID     string
}

type Post struct {
	Metadata PostMetadata `json:"metadata"`
	Filepath string
}

type PostImage struct {
	*ast.Leaf
	ImageURL string
}

var imagePrefix = []byte("![")

func parseImage(data []byte) (ast.Node, []byte, int) {
	if !bytes.HasPrefix(data, imagePrefix) {
		return nil, nil, 0
	}

	fmt.Printf("Found a gallery!\n\n")
	i := len(imagePrefix)
	// find empty line
	// TODO: should also consider end of document
	end := bytes.Index(data[i:], []byte("]"))
	if end < 0 {
		return nil, data, 0
	}
	end = end + i
	res := &PostImage{
		ImageURL: string(data[i:end]),
	}
	return res, nil, end
}

func parserHook(data []byte) (ast.Node, []byte, int) {
	if node, d, n := parseImage(data); node != nil {
		return node, d, n
	}
	return nil, nil, 0
}

func galleryRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if _, ok := node.(*PostImage); ok {
		if entering {
			// note: just for illustration purposes
			// actual implemenation of gallery in HTML / JavaScript is long
			io.WriteString(w, fmt.Sprintf(`<img src="/static/%s"></img>`, node.(*PostImage).ImageURL))
		}
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func (post *Post) RenderHTML() string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)

	contents, err := os.ReadFile(post.Filepath)

	if err != nil {
		return "<h3>Failed To Render Post</h3>"
	}

	var _meta PostMetadata

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
	var meta PostMetadata
	var post Post
	_, err := frontmatter.MustParse(strings.NewReader(md), &meta)
	if err != nil {
		return post, err
	}

	post.Metadata = meta
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
