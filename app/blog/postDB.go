package blog

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type PostDB struct {
	db *sql.DB
}

func ConnectToPostDB(filepath string) (*PostDB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	return &PostDB{
		db: db,
	}, nil
}

func (pdb *PostDB) Init() error {
	tx, err := pdb.db.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec(`CREATE TABLE IF NOT EXISTS "posts" (
	"ID"	INTEGER,
	"Title"	TEXT DEFAULT 'Empty Post',
	"Author"	TEXT DEFAULT 'George Venn',
	"Contents"	TEXT DEFAULT '# Empty Post',
	PRIMARY KEY("ID" AUTOINCREMENT)
	)
`)

	log.Println(res)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (pdb *PostDB) GetAllPosts() ([]Post, error) {
	tx, err := pdb.db.Begin()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(`SELECT * FROM "posts"`)
	if err != nil {
		return nil, err
	}

	var posts []Post

	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Metadata.ID, &post.Metadata.Title, &post.Metadata.Author, &post.Content)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)

	}

	return posts, nil

}

func (pdb *PostDB) GetPostByID(id int) ([]Post, error) {
	tx, err := pdb.db.Begin()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(`SELECT * FROM "posts" WHERE ID == ?`, id)
	if err != nil {
		return nil, err
	}

	var posts []Post

	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Metadata.ID, &post.Metadata.Title, &post.Metadata.Author, &post.Content)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)

	}

	return posts, nil

}

func (pdb *PostDB) CreateNewPost(post Post) (*sql.Result, error) {
	tx, err := pdb.db.Begin()
	if err != nil {
		return nil, err
	}

	res, err := tx.Exec(`INSERT INTO "posts" 
	(Title, Author, Contents)
	VALUES (?, ?, ?)`,
		post.Metadata.Title,
		post.Metadata.Author,
		post.Content)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &res, nil
}
