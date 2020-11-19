package router

import (
	"github.com/franciscolmos/go-meli-integration/pkg/controller"
	"github.com/gin-gonic/gin"
)

func RunAPI() {
	r := gin.Default()
	r.GET("api/auth/code", controller.GetToken)
	r.GET("api/dashboard", controller.GetDashboard)
	r.POST("api/post/item", controller.PostItem)

	r.Run( ":8080")
}

