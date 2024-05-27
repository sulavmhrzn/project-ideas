package main

import (
	"errors"
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

	err = app.writeJSON(w, http.StatusCreated, idea)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listIdeasHandler(w http.ResponseWriter, r *http.Request) {
	ideas, err := app.models.Idea.List()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, ideas)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getIdeaHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	idea, err := app.models.Idea.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRows):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, idea)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteIdeaHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Idea.Delete(id, user.Id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRows):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
