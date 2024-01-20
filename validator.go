package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type JsonErrorType struct {
	Error string `json:"error"`
}

type StatusValidType struct {
	CleanedBody string `json:"cleaned_body"`
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

func cleanBody(body string) string {
	badWords := map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}
	parts := strings.Split(body, " ")
	for i, part := range parts {
		_, ok := badWords[strings.ToLower(part)]
		if ok {
			parts[i] = "****"
		}
	}

	return strings.Join(parts, " ")
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

	cleanedBody := cleanBody(params.Body)

	if len(cleanedBody) == 0 {
		responseWithError(w, http.StatusBadRequest, "Chirp is empty")
		return
	}

	if len(cleanedBody) > 140 {
		responseWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedResponse := StatusValidType{
		CleanedBody: cleanedBody,
	}
	validResponseBody, err := json.Marshal(cleanedResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(validResponseBody)
}
