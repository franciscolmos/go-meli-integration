package controller

import (
	"github.com/franciscolmos/go-meli-integration/pkg/service"
	"github.com/gin-gonic/gin"
)

func GetDashboard (c *gin.Context){
	Dashboard := service.GetInfo( c, TokenR.Access_token, TokenR.User_id )

	c.JSON(200, Dashboard)
}