package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data any) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, input any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(input)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesReaderError *http.MaxBytesError
		switch {
		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty")
		case errors.As(err, &syntaxError):
			return fmt.Errorf("syntax error at character (%d)", syntaxError.Offset)
		case errors.As(err, &unmarshalTypeError):
			return fmt.Errorf("invalid type for field %s", unmarshalTypeError.Field)
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			unknownField := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("json contains an unknown field: %s", unknownField)
		case errors.As(err, &maxBytesReaderError):
			return fmt.Errorf("request body is too large")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	return nil
}
