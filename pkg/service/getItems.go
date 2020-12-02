package service

import (
	"encoding/json"
	"fmt"
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

// ESTRUCTURA PARA ENVIAR AL FRONT

type Item struct {
	Id string
	Title string
	Quantity int
	Price float64
	FirstPicture string
}

var itemsIds ItemsIdMeli

func GetItems( channel chan [] Item ) {

	// Obtenemos listado de ids de items del vendedor con id de vendedor y accessToken dinamicos

	ids, err := http.Get("https://api.mercadolibre.com/users/"+ strconv.Itoa(UserID) +"/items/search?access_token=" + Token)

	if err != nil {
		fmt.Errorf("Error", err.Error())
		return
	}

	defer ids.Body.Close()

	dataItemsId, err := ioutil.ReadAll(ids.Body)

	json.Unmarshal(dataItemsId, &itemsIds)

	// Listado de productos (Título, Cantidad, Precio, Primera foto)
	var items [] Item
	for i := 0; i < len(itemsIds.Id); i++ {
		resp2, err := http.Get("https://api.mercadolibre.com/items/" + itemsIds.Id[i] + "?access_token=" + Token)
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

	channel <- items
}