package router

import (
	"github.com/franciscolmos/go-meli-integration/pkg/controller"
	"github.com/gin-gonic/gin"
)

func RunAPI() {
	r := gin.Default()
	r.GET("/auth/code", controller.GetToken)
	r.GET("/dashboard", controller.GetDashboard)
	r.POST("/post/item", controller.PostItem)

	r.Run( ":8080")
}

