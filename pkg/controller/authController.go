package controller

import (
	"github.com/gin-gonic/gin"
)

func GetToken(c *gin.Context) {
	accesToken := c.Query("code")
	println( accesToken )
}
