package blog

import (
	"crypto/sha1"
	"time"
)

type Post struct {
	ID       string `json:"ID"`
	Filepath string `json:"filepath"`
	Author   string `yaml:"Author"`
	Title    string `yaml:"Title"`
}

func GetPostID(now time.Time) string {
	utc := now.UTC()

}
