package rest

import (
	"encoding/json"
	"log"
	"net/http"
)

type JSONResponse struct {
	StatusCode int `json:"status_code"`
	Message    any `json:"message"`
}

type ErrorResponse struct {
	StatusCode int `json:"status_code"`
	Error      any `json:"error"`
	Details    any `json:"details"`
}

func Send(w http.ResponseWriter, message any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.WriteHeader(statusCode)

	resp := JSONResponse{
		StatusCode: statusCode,
		Message:    message,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err)
	}
}

func Error(w http.ResponseWriter, statusCode int, errMessage string, details any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.WriteHeader(statusCode)

	resp := ErrorResponse{
		StatusCode: statusCode,
		Error:      errMessage,
		Details:    details,
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err)
	}
}
