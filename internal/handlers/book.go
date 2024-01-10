package handlers

import (
	"BookApp/book"
	"net/http"
)

type BookHandler struct {
	BookService book.UseCase
}

func NewBookHandler(service book.UseCase) *BookHandler {
	return &BookHandler{BookService: service}
}

func (b *BookHandler) FindById(w http.ResponseWriter, r *http.Request) {

}
