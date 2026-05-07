package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//go:embed _view/*
var viewsFS embed.FS

//go:embed _public/*
var publicFS embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("dotenv load failed", "err", err)
		os.Exit(1)
	}

	gin.SetMode(os.Getenv("GIN_MODE"))
	r := gin.Default()

	tmpl := template.Must(template.ParseFS(viewsFS, "_view/*.html", "_view/components/*.html"))
	r.SetHTMLTemplate(tmpl)

	publicSubFS, err := fs.Sub(publicFS, "_public")
	if err != nil {
		slog.Error("_public sub failed", "err", err)
		os.Exit(1)
	}

	tmplData := gin.H{
		"ImprintUrl":       os.Getenv("IMPRINT_URL"),
		"PrivacyPolicyUrl": os.Getenv("PRIVACY_POLICY_URL"),
	}

	links := map[string]string{}
	lock := &sync.RWMutex{}

	if err := refreshLinks(&links); err != nil {
		slog.Error("initial links refresh failed", "err", err)
		os.Exit(1)
	}

	go func() {
		for {
			if err := watchLinks(&links, lock); err != nil {
				slog.Error("watch failed", "err", err)
			}
		}
	}()

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", tmplData)
	})

	r.StaticFS("/public/", http.FS(publicSubFS))

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		lock.RLock()
		target, ok := links[path]
		lock.RUnlock()

		if ok {
			c.Redirect(http.StatusFound, target)

			return
		}

		c.HTML(http.StatusNotFound, "404.html", tmplData)
	})

	if err := r.Run(); err != nil {
		slog.Error("webserver run failed", "err", err)
		os.Exit(1)
	}
}
