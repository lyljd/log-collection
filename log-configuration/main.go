package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLFiles("./log-configuration/ui.html")
	r.Static("/static", "./log-configuration/static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "ui.html", nil)
	})

	_ = r.Run(":1295")
}
