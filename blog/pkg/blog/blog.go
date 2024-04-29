package blog

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Blog struct {
	posts    []Post
	filepath string
}

func Connect(uri string) (*Blog, error) {
	var posts []Post
	file, err := OpenOrCreate(uri)

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(data, &posts)

	return &Blog{posts: posts, filepath: uri}, nil
}

func (b *Blog) Save() {
	data, err := json.Marshal(b.posts)
	if err != nil {
		log.Println(err)
		return
	}

	os.WriteFile(b.filepath, data, os.ModePerm)
}

func (b *Blog) GetPostByID(id string) ([]Post, error) {

	var posts []Post

	for _, post := range b.posts {
		if post.ID == id {
			posts = append(posts, post)
		}
	}

	return posts, nil

}

func (b *Blog) DeletePostByID(id string) error {

	var posts []Post

	for _, post := range b.posts {
		if post.ID != id {
			posts = append(posts, post)
		}
	}

	b.posts = posts

	return nil

}

func (b *Blog) GetPostByFilePath(filepath string) ([]Post, error) {

	var posts []Post

	for _, post := range b.posts {
		if post.Filepath == filepath {
			posts = append(posts, post)
		}
	}

	return posts, nil

}

func (b *Blog) GetAllPosts() ([]Post, error) {
	return b.posts, nil
}

func (b *Blog) CreateNewPost(post Post) error {
	for i, p := range b.posts {
		if p.ID == post.ID {
			b.posts[i] = post
			return nil
		}
	}

	b.posts = append(b.posts, post)
	b.Save()
	return nil
}

type EmptyResult struct {
}

func (er *EmptyResult) LastInsertId() (int64, error) {
	return 0, fmt.Errorf("Empty Result, Check Exec's Error!\n")
}

func (er *EmptyResult) RowsAffected() (int64, error) {
	return 0, fmt.Errorf("Empty Result, Check Exec's Error!\n")
}
