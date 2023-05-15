package api

import (
	"github.com/gin-gonic/gin"
	"log-configuration/serializer"
	"log-configuration/service/key"
)

func GetKeys(c *gin.Context) {
	var s key.GetKeysService
	if err := c.ShouldBind(&s); err == nil {
		res := s.GetKeys()
		c.JSON(200, res)
	} else {
		c.JSON(200, serializer.CliParErr("", err))
	}
}
