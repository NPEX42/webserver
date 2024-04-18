package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/NPEX42/webserver/app/blog"

	blog_comps "github.com/NPEX42/webserver/components/blog"
	"github.com/NPEX42/webserver/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func CreateRouter(conf *ServerConfig, router *http.ServeMux, pdb *blog.PostDB) *http.ServeMux {

	router.Handle("GET /hooks/gh_push", middleware.Logging(logger, WebhookPushHandler()))
	router.Handle("/hooks/pull", middleware.Logging(logger, http.HandlerFunc(Restart)))
	router.Handle("/api/v1/projects", middleware.Logging(logger, http.HandlerFunc(GetProjects)))

	router.Handle("GET /blog/posts/{id}", middleware.Logging(logger, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		id := req.PathValue("id")
		log.Printf("ID: %s\n", id)
		posts, err := pdb.GetPostByID(id)
		if err != nil {
			log.Printf("ERR: %s\n", err)
			w.WriteHeader(500)
			return
		}

		fmt.Fprint(w, posts[0].RenderHTML())

	})))

	router.Handle("/", RequestLogger(logger, http.FileServer(http.Dir(conf.StaticDir))))

	return router
}

func CreateEcho(conf *ServerConfig, e *echo.Echo, pdb *blog.PostDB) {

	e.Static("/", "./static")

	e.GET("/hooks/pull", echo.WrapHandler(http.HandlerFunc(Restart)))

	e.GET("/blog/posts/:id", func(c echo.Context) error {
		id := c.Param("id")

		post, err := pdb.GetPostByID(id)
		if err != nil || len(post) == 0 {
			return c.HTML(http.StatusNotFound, fmt.Sprintf("<h1>Failed To Find Post %s</h1>", id))

		}

		var sb strings.Builder
		blog_comps.Post(post[0]).Render(context.Background(), &sb)
		return c.HTML(http.StatusOK, sb.String())
	})

	e.GET("/blog/posts", func(c echo.Context) error {
		posts, err := pdb.GetAllPosts()
		if err != nil || len(posts) == 0 {
			return c.HTML(http.StatusNotFound, fmt.Sprintf("<h1>No Posts Found</h1>"))

		}
		var sb strings.Builder
		for _, post := range posts {
			blog_comps.PostPreview(post).Render(context.Background(), &sb)
		}
		return c.HTML(http.StatusOK, sb.String())

	})

	e.GET("/api/v1/projects", echo.WrapHandler(http.HandlerFunc(GetProjects)))

}
