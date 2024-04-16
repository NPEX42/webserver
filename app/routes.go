package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/NPEX42/webserver/app/blog"
	"github.com/NPEX42/webserver/pkg/middleware"
)

func CreateRouter(conf *ServerConfig, router *http.ServeMux, pdb *blog.PostDB) *http.ServeMux {

	router.Handle("GET /hooks/gh_push", middleware.Logging(logger, WebhookPushHandler()))
	router.Handle("/hooks/pull", middleware.Logging(logger, http.HandlerFunc(Restart)))
	router.Handle("/api/v1/projects", middleware.Logging(logger, http.HandlerFunc(GetProjects)))

	router.Handle("GET /blog/posts/{id}", middleware.Logging(logger, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		id, _ := strconv.Atoi(req.PathValue("id"))
		log.Printf("ID: %d\n", id)
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
