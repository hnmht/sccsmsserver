package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"sccsmsserver/pub"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//go:embed dist/*
var embeddedFiles embed.FS

func AddRoutes(router *gin.Engine) {
	// Get "dist" file archive
	buildFS, err := fs.Sub(embeddedFiles, "dist")
	if err != nil {
		zap.L().Error("AddRoutes fs.Sub failed:", zap.Error(err))
		return
	}
	// Static resource path
	router.StaticFS("/", http.FS(buildFS))
	// Fallback
	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, pub.APIPath) {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "404 API endpoint not found",
			})
			return
		}
		file, err := buildFS.Open("index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "index.html not found")
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			zap.L().Error("AddRoutes file.Stat failed:", zap.Error(err))
			c.String(http.StatusInternalServerError, "Failed to get index.html info")
			return
		}
		c.DataFromReader(http.StatusOK, stat.Size(), "text/html", file, nil)
	})
}
