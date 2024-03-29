package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cbrgm/githubevents/githubevents"
	"github.com/google/go-github/v60/github"
	strftime "github.com/itchyny/timefmt-go"
)

var (
	Handle githubevents.EventHandler
	GH     github.Client
)

func init() {
	Handle = *githubevents.New("dd3d80f7f36a1af8ddf1cb0747051d882acebdb4c047792265f1f4f8679cc0826d64ea64f9ef8cc2e0fa93ceb7106597780895605c5e42c453878108ebe35349")
	Handle.OnPushEventAny(func(_, _ string, event *github.PushEvent) error {
		log.Printf("Repo %v Pushed.", *event.Repo.Name)
		if err := Pull(); err != nil {
			return err
		}
		return nil
	})

	Handle.OnAfterAny(func(_, _ string, _ any) error {
		time.Sleep(10 * time.Second)
		if err := Pull(); err != nil {
			return err
		}
		return nil
	})

	GH = *github.NewClient(nil)
}

func RequestLogger(logger *log.Logger, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(wrw, r)
		end := time.Now()
		duration := end.Sub(start)
		logger.Printf("%s,%s,%s,%q,%d,%d", strftime.Format(time.Now(), "%Y/%m/%d-%H:%M:%S"), r.Method, r.RemoteAddr, r.URL.Path, wrw.statusCode, duration.Milliseconds())
	}
}

func WebhookPushHandler() http.Handler {
	return http.Handler(http.HandlerFunc(WebhookPush))
}

func WebhookPush(_ http.ResponseWriter, r *http.Request) {
	err := Handle.HandleEventRequest(r)
	fmt.Println(err)
}

func GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := LoadProjects("./static/projects.json")
	if err != nil {
		w.WriteHeader(500)
		return
	}

	var s bytes.Buffer
	for _, proj := range projects {
		s.WriteString(proj.Render())
	}

	fmt.Fprintf(w, "%s", s.String())
}
