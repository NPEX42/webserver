package main

import (
	"net/http"

	"github.com/NPEX42/webserver/pkg/middleware"
)

func CreateRouter(conf *ServerConfig, router *http.ServeMux) *http.ServeMux {

	router.Handle("GET /hooks/gh_push", middleware.Logging(logger, WebhookPushHandler()))
	router.Handle("/hooks/pull", middleware.Logging(logger, http.HandlerFunc(Restart)))
	router.Handle("/api/v1/projects", middleware.Logging(logger, http.HandlerFunc(GetProjects)))

	router.Handle("/", RequestLogger(logger, http.FileServer(http.Dir(conf.StaticDir))))

	return router
}
