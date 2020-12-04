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

		var UnansweredQuestiontemp Unanswered_Question

		for i := 0; i < len(questions.Questions); i++ {
			UnansweredQuestiontemp.Id = questions.Questions[i].Id
			if len(questions.Questions) == 0 || questions.Questions[i].Status != "UNANSWERED" {
				continue
			}
			for j := 0; j < len(Dash.Items); j++ {
				if Dash.Items[j].Id == questions.Questions[i].Item_id {
					UnansweredQuestiontemp.Title = Dash.Items[j].Title
				}
			}
			UnansweredQuestiontemp.Question_date = questions.Questions[i].Date_created
			UnansweredQuestiontemp.Question_text = questions.Questions[i].Text

			Unanswered_Questions = append(Unanswered_Questions, UnansweredQuestiontemp)

			question := model.Question{
				Text: UnansweredQuestiontemp.Question_text,
				Question_Id: UnansweredQuestiontemp.Id,
				ItemTitle: UnansweredQuestiontemp.Title,
				CreatedAt:time.Now(),
				UpdatedAt: time.Now() }

			var questions [] model.Question


			//Consultamos si existe un item con el id que nos devuelve meli
			db.Where("question_id = ?",UnansweredQuestiontemp.Id).First(&questions)

			//en caso de exista, entonces continuamos con el que sigue, pero si no existe, lo agregamos a la base de datos.
			if len(questions) == 0 {
				db.Create(&question)
			}
		}
	}

	return Unanswered_Questions

//	channel <- Unanswered_Questions
//	close(channel)
}