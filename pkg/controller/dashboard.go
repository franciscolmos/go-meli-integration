package controller
import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

// ITEMS DEL VENDEDOR
type ItemsIdMeli struct {
	Id []string              `json:"results"`
}

type PictureMeli struct {
	Url string                 `json:"url"`
}

type ItemMeli struct {
	Id    string               `json:"id"`
	Title string               `json:"title"`
	Price float64              `json:"price"`
	Available_quantity int     `json:"available_quantity"`
	Pictures []PictureMeli	   `json:"pictures"`
}

// ITEMS VENDIDOS
type SingleItemMeli struct {
	Title string                `json:"title"`
}

type Order_ItemsMeli struct {
	SingleItem SingleItemMeli    `json:"item"`
	Quantity int                 `json:"quantity"`
	Unit_price float64           `json:"unit_price"`
	Full_Unit_Price float64      `json:"full_unit_price"`
}

type ResultMeli struct {
	Order_Items []Order_ItemsMeli `json:"order_items"`
	Total_amount float64          `json:"total_amount"`
	Paid_amount float64           `json:"paid_amount"`
	Date_closed string            `json:"date_closed"`
}

type SoldItemMeli struct {
	Result []ResultMeli            `json:"results"`
}

// PREGUNTAS SIN RESPONDER
type QuestionMeli struct {
	Item_id string        `json:"item_id"`
	Date_created string   `json:"date_created"`
	Text string           `json:"text"`
	Status string         `json:"status"`
}

type QuestionsMeli struct {
	Questions []QuestionMeli  `json:"questions"`
}

// ESTRUCTURA PARA ENVIAR AL FRONT

type Item struct {
	Id string
	Title string
	Quantity int
	Price float64
	FirstPicture string
}

type Sold_Item struct {
	Title string
	Sold_Quantity int
	Unit_Price float64
	Subtotal float64
}

type Sale_Order struct {
	Sold_Items [] Sold_Item
	Sale_date string
	Total  float64
	Total_Delivery float64
}

type Unanswered_Question struct {
	Question_date string
	Title string
	Question_text string
}

type Dashboard struct {
	Items [] Item
	Sales_Orders [] Sale_Order
	Unanswered_Questions [] Unanswered_Question
}

func GetDashboard (c *gin.Context){
	var Dashboard Dashboard

	// Obtenemos listado de ids de items del vendedor con id de vendedor y accessToken dinamicos

	ids, err := http.Get("https://api.mercadolibre.com/users/"+ strconv.Itoa(TokenR.User_id) +"/items/search?access_token=" + TokenR.Access_token)

	if err != nil {
		fmt.Errorf("Error", err.Error())
		return
	}

	defer ids.Body.Close()

	dataItemsId, err := ioutil.ReadAll(ids.Body)

	var itemsIds ItemsIdMeli
	json.Unmarshal(dataItemsId, &itemsIds)

	// Listado de productos (Título, Cantidad, Precio, Primera foto)
	var items [] Item
	for i := 0; i < len(itemsIds.Id); i++ {
		resp2, err := http.Get("https://api.mercadolibre.com/items/" + itemsIds.Id[i] + "?access_token=" + TokenR.Access_token)
		if err != nil {
			fmt.Errorf("Error",err.Error())
			return
		}
		dataItemsDetail, err := ioutil.ReadAll(resp2.Body)

		var item ItemMeli
		json.Unmarshal(dataItemsDetail, &item)

		var itemTemp Item

		itemTemp.Id = item.Id
		itemTemp.Title = item.Title
		itemTemp.Price = item.Price
		itemTemp.FirstPicture = item.Pictures[0].Url
		itemTemp.Quantity = item.Available_quantity

		items = append(items, itemTemp)
	}

	//asignamos al dashboard los items
	Dashboard.Items = items

	//  Ventas efectuadas
	resp2, err := http.Get("https://api.mercadolibre.com/orders/search?seller="+ strconv.Itoa(TokenR.User_id) +"&order.status=paid&access_token=" + TokenR.Access_token)

	defer resp2.Body.Close()

	dataSales, err := ioutil.ReadAll(resp2.Body)

	var soldItems SoldItemMeli
	json.Unmarshal(dataSales, &soldItems)

	var Sales_Orders [] Sale_Order

	for i := 0; i < len(soldItems.Result); i++ {
		var Sale_Order_Temp Sale_Order
		Sale_Order_Temp.Sale_date = soldItems.Result[i].Date_closed
		Sale_Order_Temp.Total = soldItems.Result[i].Total_amount
		Sale_Order_Temp.Total_Delivery = soldItems.Result[i].Paid_amount
		for j := 0; j < len(soldItems.Result[i].Order_Items); j++ {
			var Sale_Order_Temp_Items Sold_Item
			Sale_Order_Temp_Items.Title = soldItems.Result[i].Order_Items[j].SingleItem.Title
			Sale_Order_Temp_Items.Unit_Price = soldItems.Result[i].Order_Items[j].Unit_price
			Sale_Order_Temp_Items.Sold_Quantity = soldItems.Result[i].Order_Items[j].Quantity
			Sale_Order_Temp_Items.Subtotal = soldItems.Result[i].Order_Items[j].Full_Unit_Price

			Sale_Order_Temp.Sold_Items = append(Sale_Order_Temp.Sold_Items, Sale_Order_Temp_Items)
		}
		Sales_Orders = append(Sales_Orders, Sale_Order_Temp)
	}

	Dashboard.Sales_Orders = Sales_Orders

	// Preguntas pendientes por responder por cada ítem ordenadas de las más antiguas a las más recientes.
	var Unanswered_Questions []Unanswered_Question

	for i := 0; i < len(itemsIds.Id); i++ {
		resp3, err := http.Get("https://api.mercadolibre.com/questions/search?item=" + itemsIds.Id[i] + "&access_token=" + TokenR.Access_token + "&sort_fields=date_created&sort_types=ASC")
		if err != nil {
			fmt.Errorf("Error", err.Error())
			return
		}
		dataQuestions, err := ioutil.ReadAll(resp3.Body)

		var questions QuestionsMeli
		json.Unmarshal(dataQuestions, &questions)

		var UnansweredQuestiontemp Unanswered_Question

		for i := 0; i < len(questions.Questions); i++ {
			if  len(questions.Questions) == 0 || questions.Questions[i].Status != "UNANSWERED" {
				continue
			}
			for j := 0; j < len(Dashboard.Items); j++ {
				if Dashboard.Items[j].Id == questions.Questions[i].Item_id {
					UnansweredQuestiontemp.Title = Dashboard.Items[j].Title
				}
			}
			UnansweredQuestiontemp.Question_date = questions.Questions[i].Date_created
			UnansweredQuestiontemp.Question_text = questions.Questions[i].Text

			Unanswered_Questions = append(Unanswered_Questions, UnansweredQuestiontemp)
		}
	}

	Dashboard.Unanswered_Questions = Unanswered_Questions

	c.JSON(200, Dashboard)
}
