package blog

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type PostDB struct {
	db *sql.DB
}

func ConnectToPostDB(filepath string) (*PostDB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	log.Println("Opened DB")

	return &PostDB{
		db: db,
	}, nil
}

func (pdb *PostDB) Init() error {

	_, err := pdb.db.Exec(`CREATE TABLE IF NOT EXISTS "posts" (
	"ID"		TEXT,
	"Title"		TEXT DEFAULT 'Empty Post',
	"Author"	TEXT DEFAULT 'George Venn',
	"Filepath"	TEXT DEFAULT '',
	"Views"		INTEGER DEFAULT 0,
	PRIMARY KEY("ID")
	)
	`)

	if err != nil {
		return err
	}

	return nil
}

func (pdb *PostDB) GetAllPosts() ([]Post, error) {

	rows, err := pdb.db.Query(`SELECT * FROM "posts"`)
	if err != nil {
		return nil, err
	}

	var posts []Post

	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Metadata.ID, &post.Metadata.Title, &post.Metadata.Author, &post.Filepath)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)

	}

	return posts, nil

}

func (pdb *PostDB) GetPostByID(id string) ([]Post, error) {

	rows, err := pdb.db.Query(`SELECT * FROM "posts" WHERE ID == ?`, id)
	if err != nil {
		return nil, err
	}

	var posts []Post

	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Metadata.ID, &post.Metadata.Title, &post.Metadata.Author, &post.Filepath)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)

	}

	return posts, nil

}

func (pdb *PostDB) GetPostByTitle(title string) ([]Post, error) {

	rows, err := pdb.db.Query(`SELECT * FROM "posts" WHERE Title == ?`, title)
	if err != nil {
		return nil, err
	}

	var posts []Post

	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Metadata.ID, &post.Metadata.Title, &post.Metadata.Author, &post.Filepath)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)

	}

	return posts, nil

}

func (pdb *PostDB) GetPostByFilepath(title string) ([]Post, error) {

	rows, err := pdb.db.Query(`SELECT * FROM "posts" WHERE Filepath == ?`, title)
	if err != nil {
		return nil, err
	}

	var posts []Post

	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Metadata.ID, &post.Metadata.Title, &post.Metadata.Author, &post.Filepath)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)

	}

	return posts, nil

}

func (pdb *PostDB) DeletePostByID(id string) error {
	_, err := pdb.db.Exec("DELETE FROM posts WHERE ID == ?", id)
	return err
}

func (pdb *PostDB) CreateNewPost(post Post) (*sql.Result, error) {

	return nil, sql.ErrNoRows
}

func (pdb *PostDB) Close() {
	log.Printf("Closing DB")
	pdb.db.Close()
}
