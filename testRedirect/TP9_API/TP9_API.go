package TP9_API

import (
	"github.com/franciscolmos/go-meli-integration/testRedirect/TP9_API/controller"
	"github.com/gin-gonic/gin"
)

func RunAPI() {
	r := gin.Default()
	r.GET("/auth/code", controller.GetToken)
	r.Run( ":8080")
}
