package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NPEX42/webserver/app/blog"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Static("/", "./static")

	pdb, err := blog.ConnectToPostDB("./posts.db")
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	e.GET("/blog/posts/:id", func(c echo.Context) error {
		id := c.Param("id")

		post, err := pdb.GetPostByID(id)
		if err != nil || len(post) == 0 {
			return c.HTML(http.StatusNotFound, fmt.Sprintf("<h1>Failed To Find Post %s</h1>", id))

		}

		return c.HTML(http.StatusOK, post[0].RenderHTML())
	})

	e.Start(":8080")
}
