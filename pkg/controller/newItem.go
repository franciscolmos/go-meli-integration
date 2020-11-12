package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

//SUBESTRUCTURAS NECESARIAS
type Description struct {
	Plain_text string `json:"plain_text"`
}

type SaleTerms struct {
	Id string `json:"id"`
	Value_name string `json:"value_name"`
}

type Picture struct {
	Source string `json:"source"`
}

type Atribute struct {
	Id string `json:"id"`
	Value_name string `json:"value_name"`
}

//ESTRUCTURA QUE SE MANDA POR POST
type NewItem struct {
	Title string `json:"title"`
	CategoryId string `json:"category_id"`
	Price float64 `json:"price"`
	Currency_id string `json:"currency_id"`
	Available_quantity int `json:"available_quantity"`
	Buying_mode string `json:"buying_mode"`
	Condition string `json:"condition"`
	Listing_type_id string `json:"listing_type_id"`
	Description Description `json:"description"`
	Video_id string `json:"video_id"`
	Sale_terms [] SaleTerms `json:"sale_terms"`
	Pictures [] Picture `json:"pictures"`
	Atributes [] Atribute `json:"attributes"`
}

//ESTRUCTURA DE RESPUESTA
type Response struct {
	Id string `json:"id"`
	Title string `json:"title"`
}

var ResponseNewItem Response

func PostItem ( c *gin.Context) {
	atributes := []Atribute{
		{
			Id: "BRAND",
			Value_name: "Marca del producto",
		},
		{
			Id: "EAN",
			Value_name: "7898095297749",
		},
	}

	pictures := []Picture {
		{
			Source: "http://mla-s2-p.mlstatic.com/968521-MLA20805195516_072016-O.jpg",
		},
	}

	sale_terms := []SaleTerms{
		{
			Id: "WARRANTY_TYPE",
			Value_name: "Garantía del vendedor",
		},
		{
			Id: "WARRANTY_TIME",
			Value_name: "90 días",
		},
	}

	description := Description {
		Plain_text: "Descripción con Texto Plano \n",
	}

	newItem := NewItem{
		Title: "item de prueba",
		CategoryId: "MLA3530",
		Price: 2960,
		Currency_id: "ARS",
		Available_quantity: 20,
		Buying_mode: "buy_it_now",
		Condition: "new",
		Listing_type_id: "gold_special",
		Description: description,
		Video_id: "YOUTUBE_ID_HERE",
		Sale_terms: sale_terms,
		Pictures: pictures,
		Atributes: atributes,
	}

	jsonNewItem, _ := json.Marshal(newItem)

	fmt.Println(string(jsonNewItem))

	responsePostNewItem, err := http.Post("https://api.mercadolibre.com/items?access_token=" + TokenR.Access_token, "application/json; application/x-www-form-urlencoded", bytes.NewBuffer(jsonNewItem))

	if err != nil {
		fmt.Errorf("Error", err.Error())
		return
	}

	defer responsePostNewItem.Body.Close()

	data, err := ioutil.ReadAll(responsePostNewItem.Body)

	if err != nil {
		fmt.Errorf("Error", err.Error())
		return
	}

	bodyString := string(data)
	fmt.Println(bodyString)

	json.Unmarshal(data, &ResponseNewItem)

	c.JSON(200, ResponseNewItem)

}