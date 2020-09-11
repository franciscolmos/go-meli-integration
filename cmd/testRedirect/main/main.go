package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"github.com/gin-gonic/gin"
)

func redirect(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "https://auth.mercadolibre.com.ar/authorization?response_type=code&" +
		"client_id=2760149476611182&redirect_uri=http://localhost:8080/", 301)
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func getToken(c *gin.Context) {
	accesToken := c.Query("code")
	println( "123" )
	panic(accesToken)
}

func main() {
	r := gin.Default()
	http.HandleFunc("/", redirect)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	openbrowser("https://auth.mercadolibre.com.ar/authorization?response_type=code&" +
		"client_id=2760149476611182&redirect_uri=http://localhost:8080/")
	r.GET("http://localhost:8080/", getToken )
	print( getToken )
}