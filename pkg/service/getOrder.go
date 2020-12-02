package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

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

// ESTRUCTURA PARA ENVIAR AL FRONT
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

func GetOrder( channel chan [] Sale_Order ) {
	//  Ventas efectuadas
	resp2, err := http.Get("https://api.mercadolibre.com/orders/search?seller="+ strconv.Itoa(UserID) +"&order.status=paid&access_token=" + Token)

	if err != nil {
		fmt.Errorf("Error", err.Error())
		return
	}

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

	channel <- Sales_Orders
}
