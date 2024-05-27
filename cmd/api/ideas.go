package main

import (
	"fmt"
	"net/http"

	"github.com/sulavmhrzn/projectideas/internal/data"
	"github.com/sulavmhrzn/projectideas/internal/validator"
)

func (app *application) createIdeaHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	var input struct {
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Tags        []data.Tag `json:"tags"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	idea := &data.Idea{
		Title:       input.Title,
		Description: input.Description,
		Tags:        input.Tags,
		UserId:      user.Id,
	}
	if data.ValidateIdea(v, idea); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	idea, err = app.models.Idea.Insert(idea)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	fmt.Printf("%+v", idea)
	w.Write([]byte(fmt.Sprintf("%+v", idea)))

}
