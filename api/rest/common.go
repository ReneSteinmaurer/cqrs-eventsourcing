package rest

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response[T any] struct {
	Errors  []string `json:"errors"`
	Data    T        `json:"data"`
	Message string   `json:"message"`
}

func NewResponseContent[T any](data T) Response[T] {
	return Response[T]{
		Errors: make([]string, 0),
		Data:   data,
	}
}

func NewResponseContentMessage[T any](msg string, data T) Response[T] {
	return Response[T]{
		Errors:  make([]string, 0),
		Data:    data,
		Message: msg,
	}
}

func NewResponseErrors(errors ...string) Response[any] {
	return Response[any]{
		Errors: errors,
	}
}

func SendResponseErrors(w *http.ResponseWriter, errors ...string) {
	if w == nil {
		log.Println("[ERROR] SendResponse: ResponseWriter is nil")
		return
	}
	res := NewResponseErrors(errors...)
	SendResponse(res, w)
}

func SendResponseContent[T any](w *http.ResponseWriter, data T) {
	if w == nil {
		log.Println("[ERROR] SendResponse: ResponseWriter is nil")
		return
	}
	res := NewResponseContent(data)
	SendResponse(res, w)
}

func SendResponseContentMessage[T any](w *http.ResponseWriter, data T, message string) {
	if w == nil {
		log.Println("[ERROR] SendResponse: ResponseWriter is nil")
		return
	}
	res := NewResponseContentMessage(message, data)
	SendResponse(res, w)
}

func SendResponse[T any](res Response[T], w *http.ResponseWriter) {
	if w == nil {
		log.Println("[ERROR] SendResponse: ResponseWriter is nil")
		return
	}

	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusOK)

	if err := json.NewEncoder(*w).Encode(res); err != nil {
		log.Printf("[ERROR] SendResponse: Failed to encode response: %v\n", err)
	}
}
