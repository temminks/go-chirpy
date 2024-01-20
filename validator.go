package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type JsonErrorType struct {
	Error string `json:"error"`
}

type StatusValidType struct {
	Valid bool `json:"valid"`
}

func createError(text string) (jsonError []byte, err error) {
	jsonErrorType := JsonErrorType{
		Error: text,
	}

	jsonError, err = json.Marshal(jsonErrorType)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return nil, err
	}

	return jsonError, nil
}

func responseWithError(w http.ResponseWriter, code int, msg string) {
	errorResponse, err := createError(msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(errorResponse)
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := chirp{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		responseWithError(w, http.StatusInternalServerError, "Error decoding parameters.")
		return
	}

	if len(params.Body) == 0 {
		responseWithError(w, http.StatusBadRequest, "Chirp is empty")
		return
	}

	if len(params.Body) > 140 {
		responseWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	valid := StatusValidType{
		Valid: true,
	}
	validResponseBody, err := json.Marshal(valid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(validResponseBody)
}
