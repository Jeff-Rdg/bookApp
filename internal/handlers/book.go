package handlers

import (
	"BookApp/book"
	"BookApp/pkg/filter"
	"BookApp/pkg/pagination"
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
	param := chi.URLParam(r, "id")
	if param == "" {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "No id informed"})
		return
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "error to find book", Data: err})
		return
	}

	result, err := b.BookService.GetById(uint(id))
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "error to find book", Data: err})
		return
	}

	RenderJSON(w, Result{StatusCode: http.StatusOK, Data: result})
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

	RenderJSON(w, Result{StatusCode: http.StatusCreated, Message: fmt.Sprintf("book with id %v created successfully", id)})
}

func (b *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	var bookRequest book.Request
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

func (b *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	if param == "" {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "No id informed"})
		return
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "invalid id informed", Data: err})
		return
	}

	err = b.BookService.Delete(uint(id))
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "error to delete book", Data: err})
	}

	RenderJSON(w, Result{StatusCode: http.StatusOK, Message: fmt.Sprintf("book with id %v deleted successfully", id)})
}

func (b *BookHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	sort := r.URL.Query().Get("sort")
	paramAuthorName := r.URL.Query().Get("author_name")
	paramPubYear := r.URL.Query().Get("publication_year")
	paramEdition := r.URL.Query().Get("edition")
	paramName := r.URL.Query().Get("name")

	if paramPubYear == "" {
		paramPubYear = "0"
	}

	paramIntPubYear, err := strconv.Atoi(paramPubYear)
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "error to parse queries", Data: err})
		return
	}

	aux := filter.Book{
		Name:            paramName,
		Edition:         paramEdition,
		PublicationYear: uint(paramIntPubYear),
		AuthorName:      paramAuthorName,
	}

	pag, err := pagination.NewPagination(limit, page, sort)
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "error to parse queries", Data: err})
		return
	}

	result, err := b.BookService.List(*pag, aux)
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "error to list authors", Data: err})
		return
	}

	RenderJSON(w, Result{StatusCode: http.StatusOK, Data: result})
}
