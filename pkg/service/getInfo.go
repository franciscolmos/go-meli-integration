package service

import "github.com/gin-gonic/gin"

type Dashboard struct {
	Items [] Item
	Sales_Orders [] Sale_Order
	Unanswered_Questions [] Unanswered_Question
}

var Token string
var UserID int

var Dash Dashboard

func GetInfo( c *gin.Context, token string, userID int ) Dashboard {
	Token = token
	UserID = userID

	c1 := make( chan [] Item )
	c2 := make( chan [] Sale_Order )
	//c3 := make( chan [] Unanswered_Question )

	go GetItems( c1 )
	go GetOrder( c2 )
	//go getQuestion( c3 )

	Dash.Items = <- c1
	Dash.Sales_Orders = <- c2
	Dash.Unanswered_Questions = getQuestion()

	/*
	go func() {
		for i := 0; i < 20; i++ {
			select {
			case Items := <-c1:
				println( "ACA ESTAMOS EN EL CASE ITEMS")
				Dash.Items = Items
				Dash.Unanswered_Questions = getQuestion()
			case Sales_Orders := <-c2:
				println( "ACA ESTAMOS EN EL CASE SALE ORDERS")
				Dash.Sales_Orders = Sales_Orders
	//		case Unanswered_Question := <-c3:
	//			Dash.Unanswered_Questions = Unanswered_Question
			}
		}
	}()
	*/

	return Dash
}
