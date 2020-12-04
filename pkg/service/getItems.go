package service

import (
	"encoding/json"
	"fmt"
	"github.com/franciscolmos/go-meli-integration/pkg/database"
	"github.com/franciscolmos/go-meli-integration/pkg/database/model"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
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

	db := database.ConnectDB()

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

		//INSERTAMOS EN LA BASE DE DATOS LOS ITEMS QUE NO ESTEN CARGADOS
		itemDb := model.Item{ Title: item.Title,
							  Quantity: item.Available_quantity,
							  Price: item.Price,
							  FirstPicture: item.Pictures[0].Url,
							  ItemId: item.Id,
							  CreatedAt:time.Now(),
							  UpdatedAt: time.Now() }

		var items [] model.Item

		//Consultamos si existe un item con el id que nos devuelve meli
		db.Where("item_id = ?", item.Id).First(&items)

		//en caso de exista, entonces continuamos con el que sigue, pero si no existe, lo agregamos a la base de datos.
		if len(items) != 0 {
			continue
		} else{
			db.Create(&itemDb)
		}
	}

	channel <- items
}