package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

var code string

type Token struct {
	Grant_type string 		`json:"grant_type"`
	Client_id int 			`json:"client_id"`
	Client_secret string 	`json:"client_secret"`
	Code string				`json:"code"`
	Redirect_uri string 	`json:"redirect_uri"`
}

type TokenResp struct {
	Access_token string 	`json:"access_token"`
	Token_type string 		`json:"token_type"`
	Expires_in int 			`json:"expires_in"`
	Scope string 			`json:"scope"`
	User_id int 			`json:"user_id"`
	Refresh_token string 	`json:"refresh_token"`
}

type ItemsId struct {
	Id []string              `json:"results"`
}

type Picture struct {
	Url string                 `json:"url"`
}

type Item struct {
	Title string               `json:"title"`
	Price float64              `json:"price"`
	Available_quantity int     `json:"available_quantity"`
	Pictures []Picture		   `json:"pictures"`
}



func GetToken(c *gin.Context) {
	code = c.Query("code")
	fmt.Println("code: " + code)
	TokenRequest(code)
}


func TokenRequest(code string) {

	u := Token{
		Grant_type: "authorization_code",
		Client_id: 2760149476611182,
		Client_secret: "G0vTscPHYNlLrB148wwdsjuwkqWU1HOy",
		Code: code,
		Redirect_uri: "http://localhost:8080/auth/code/",
	}

	b, err := json.Marshal(u)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

	// Intercambiamos code por token
	resp, err := http.Post("https://api.mercadolibre.com/oauth/token","application/json; application/x-www-form-urlencoded", bytes.NewBuffer(b))

	if err != nil {
		fmt.Errorf("Error",err.Error())
		return
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	bodyString := string(data)
	fmt.Println(bodyString)

	var tokenResp TokenResp
	json.Unmarshal(data, &tokenResp)
	fmt.Printf("%+v\n", tokenResp)

	// Obtenemos listado de ids de items del vendedor con id de vendedor y accessToken dinamicos
	resp1, err := http.Get("https://api.mercadolibre.com/users/"+ strconv.Itoa(tokenResp.User_id) +"/items/search?access_token=" + tokenResp.Access_token)

	defer resp1.Body.Close()

	data1, err := ioutil.ReadAll(resp1.Body)

	var itemsIds ItemsId
	json.Unmarshal(data1, &itemsIds)
	fmt.Printf("%+v\n", itemsIds)


	// mostramos todos los items del vendedor
	for i := 0; i < len(itemsIds.Id); i++ {
		resp2, err := http.Get("https://api.mercadolibre.com/items/" + itemsIds.Id[i] + "?access_token=APP_USR-2760149476611182-110500-2af3caa0dd8bce845f0493a720b4d10d-651268893")
		if err != nil {
			fmt.Errorf("Error",err.Error())
			return
		}
		data2, err := ioutil.ReadAll(resp2.Body)

		var item Item
		json.Unmarshal(data2, &item)
		fmt.Printf("%+v\n", item)
	}
}
