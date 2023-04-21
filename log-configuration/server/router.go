package server

import (
	"github.com/gin-gonic/gin"
	"log-collection/log-configuration/api"
)

func Run() {
	go func() {
		r := gin.Default()

		r.LoadHTMLFiles("./log-configuration/ui.html")
		r.Static("/static", "./log-configuration/static")

		r.GET("/", func(c *gin.Context) {
			c.HTML(200, "ui.html", nil)
		})

		r.GET("/keys", api.GetKeys)

		r.GET("/configuration/:key", api.GetConfigurationByKey)

		r.PUT("/configuration/:key", api.SetConfigurationByKey)

		_ = r.Run(":1295")
	}()
}