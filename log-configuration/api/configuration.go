package api

import (
	"github.com/gin-gonic/gin"
	"log-collection/log-configuration/serializer"
	"log-collection/log-configuration/service/configuration"
)

func GetConfigurationByKey(c *gin.Context) {
	var s configuration.GetConfigurationService
	if err := c.ShouldBind(&s); err == nil {
		res := s.GetConfigurationByKey(c.Param("key"))
		c.JSON(200, res)
	} else {
		c.JSON(200, serializer.CliParErr("", err))
	}
}

func SetConfigurationByKey(c *gin.Context) {
	var s configuration.SetConfigurationService
	if err := c.ShouldBind(&s); err == nil {
		res := s.SetConfigurationByKey(c.Param("key"))
		c.JSON(200, res)
	} else {
		c.JSON(200, serializer.CliParErr("", err))
	}
}
