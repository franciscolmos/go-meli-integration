package service

import (
	"encoding/json"
	"fmt"
	"github.com/franciscolmos/go-meli-integration/pkg/database"
	"github.com/franciscolmos/go-meli-integration/pkg/database/model"
	"io/ioutil"
	"net/http"
	"time"
)

// PREGUNTAS SIN RESPONDER
type QuestionMeli struct {
	Id int                `json:"id"`
	Item_id string        `json:"item_id"`
	Date_created string   `json:"date_created"`
	Text string           `json:"text"`
	Status string         `json:"status"`
}

type QuestionsMeli struct {
	Questions []QuestionMeli  `json:"questions"`
}

// ESTRUCTURA PARA ENVIAR AL FRONT

type Unanswered_Question struct {
	Id int
	Question_date string
	Title string
	Question_text string
}

func getQuestion() []Unanswered_Question {
	// Preguntas pendientes por responder por cada ítem ordenadas de las más antiguas a las más recientes.
	var Unanswered_Questions []Unanswered_Question

	db := database.ConnectDB()

	for i := 0; i < len(itemsIds.Id); i++ {
		resp3, err := http.Get("https://api.mercadolibre.com/questions/search?item=" + itemsIds.Id[i] + "&access_token=" + Token + "&sort_fields=date_created&sort_types=ASC")
		if err != nil {
			fmt.Errorf("Error", err.Error())
			return nil
		}
		dataQuestions, err := ioutil.ReadAll(resp3.Body)

		var questions QuestionsMeli
		json.Unmarshal(dataQuestions, &questions)
		fmt.Println("Estructura questions total: ",questions)

		var UnansweredQuestiontemp Unanswered_Question

		for j:= 0; j < len(questions.Questions); j++ {
			println("Id preguntas: ",questions.Questions[j].Id)
			UnansweredQuestiontemp.Id = questions.Questions[j].Id
			if len(questions.Questions) != 0 && questions.Questions[j].Status == "UNANSWERED" {
				for k := 0; k < len(Dash.Items); k++ {
					fmt.Println("ITEMS: ",Dash.Items[k])
					if Dash.Items[k].Id == questions.Questions[j].Item_id {
						UnansweredQuestiontemp.Title = Dash.Items[k].Title
					}
				}
				UnansweredQuestiontemp.Question_date = questions.Questions[j].Date_created
				UnansweredQuestiontemp.Question_text = questions.Questions[j].Text

				fmt.Println( UnansweredQuestiontemp.Id )
				Unanswered_Questions = append(Unanswered_Questions, UnansweredQuestiontemp)

				question := model.Question{
					Text: UnansweredQuestiontemp.Question_text,
					ItemTitle: UnansweredQuestiontemp.Title,
					CreatedAt:time.Now(),
					UpdatedAt: time.Now() }

				var questions [] model.Question


				//Consultamos si existe un item con el id que nos devuelve meli
				db.Where("text = ?",UnansweredQuestiontemp.Question_text).First(&questions)

				//en caso de exista, entonces continuamos con el que sigue, pero si no existe, lo agregamos a la base de datos.
				if len(questions) == 0 {
					db.Create(&question)
				}

				UnansweredQuestiontemp.Id = 0
			}
		}
	}

	return Unanswered_Questions

//	channel <- Unanswered_Questions
//	close(channel)
}