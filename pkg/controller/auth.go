package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/franciscolmos/go-meli-integration/pkg/database"
	"github.com/franciscolmos/go-meli-integration/pkg/database/model"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

var code string
var TokenR TokenResp

// TOKEN
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

// FUNCIONES PARA INTERCAMBIAR EL CODE POR UN ACCESS TOKEN
func GetToken(c *gin.Context) {
	code = c.Query("code")
	TokenRequest(code, c)
}

func TokenRequest(code string, c *gin.Context) {
	token := Token{
		Grant_type:    "authorization_code",
		Client_id:     2760149476611182,
		Client_secret: "G0vTscPHYNlLrB148wwdsjuwkqWU1HOy",
		Code:          code,
		Redirect_uri:  "http://localhost:4200/dashboard/",
	}

	jsonToken, err := json.Marshal(token)

	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(string(jsonToken))

	// Intercambiamos code por token
	resp, err := http.Post("https://api.mercadolibre.com/oauth/token", "application/json; application/x-www-form-urlencoded", bytes.NewBuffer(jsonToken))

	if err != nil {
		fmt.Errorf("Error", err.Error())
		return
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Errorf("Error", err.Error())
		return
	}

	//bodyString := string(data)
	//fmt.Println(bodyString)

	json.Unmarshal(data, &TokenR)
	//fmt.Printf("%+v\n", TokenR)


	// REGISTRAMOS EN LA BASE DE DATOS EL USUARIO QUE ACABA DE INICIAR SESION, SI YA ESTA REGISTRADO SOLO SE ACTUALIZAN DATOS
	// EN CASO CONTRARIO SE CREA UNA NUEVA FILA.
	user := model.User{ AccessToken: TokenR.Access_token,
						RefreshToken: TokenR.Refresh_token,
						UserIdMeli: TokenR.User_id,
						CreatedAt:time.Now(),
						UpdatedAt: time.Now() }

	db := database.ConnectDB()

	var users [] model.User

	//Consultamos si ese usuario ya esta registrado en la db
	db.Where("user_id_meli = ?", TokenR.User_id).First(&users)

	//en caso de que si esta registrado, entonces se actualizan sus datos, caso contrario se crea un nuevo registro.
	if len(users) != 0 {
		db.Model(model.User{}).Where("user_id_meli = ?", TokenR.User_id).Updates(user)
	}else {
		db.Create(&user)
	}

	c.JSON(200, TokenR)
}
