package handlers

import (
	"BookApp/book"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

type BookHandler struct {
	BookService book.UseCase
}

func NewBookHandler(service book.UseCase) *BookHandler {
	return &BookHandler{BookService: service}
}

func (b *BookHandler) FindById(w http.ResponseWriter, r *http.Request) {

}

func (b *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newBook book.Request

	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "Error to decode body", Data: err})
		return
	}

	id, err := b.BookService.Create(newBook)
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "Error to create book", Data: err})
		return
	}

	RenderJSON(w, Result{StatusCode: http.StatusCreated, Message: fmt.Sprintf("book with id %v created sucessfully", id)})
}

func (b *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	var bookRequest book.UpdateRequest
	param := chi.URLParam(r, "id")
	if param == "" {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "No id informed"})
		return
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "invalid format to id", Data: err})
		return
	}

	err = json.NewDecoder(r.Body).Decode(&bookRequest)
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "Error to decode body", Data: err})
		return
	}

	updatedId, err := b.BookService.Update(uint(id), bookRequest)
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "Error to update book", Data: err})
		return
	}

	RenderJSON(w, Result{StatusCode: http.StatusOK, Message: fmt.Sprintf("book with id %v updated sucessfully", updatedId)})
}
