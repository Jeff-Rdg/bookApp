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
		NewErrorResponse("error to find book by id",
			"the requested endpoint requires the query id parameter",
			http.StatusBadRequest,
			nil, r.URL.Path).
			RenderJSON(w)
		return
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		NewErrorResponse("error to find book by id",
			"parameter with invalid format",
			http.StatusBadRequest,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}

	result, err := b.BookService.GetById(uint(id))
	if err != nil {
		NewErrorResponse("error to find book by id",
			"",
			http.StatusNotFound,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}
	NewSuccessResponse(http.StatusOK, "", result).RenderJSON(w)
}

func (b *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newBook book.Request

	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		NewErrorResponse("error to create new book",
			"error to decode body",
			http.StatusBadRequest,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}

	id, err := b.BookService.Create(newBook)
	if err != nil {
		NewErrorResponse("error to create new book",
			"",
			http.StatusBadRequest,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}

	NewSuccessResponse(http.StatusCreated, fmt.Sprintf("book with id %v created successfully", id), nil).RenderJSON(w)
}

func (b *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	var bookRequest book.Request
	param := chi.URLParam(r, "id")
	if param == "" {
		NewErrorResponse("error to update book",
			"the requested endpoint requires the query id parameter",
			http.StatusBadRequest,
			nil, r.URL.Path).
			RenderJSON(w)
		return
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		NewErrorResponse("error to update book",
			"parameter with invalid format",
			http.StatusBadRequest,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&bookRequest)
	if err != nil {
		NewErrorResponse("error to update book",
			"Error to decode body",
			http.StatusBadRequest,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}

	updatedId, err := b.BookService.Update(uint(id), bookRequest)
	if err != nil {
		NewErrorResponse("error to update book",
			"",
			http.StatusBadRequest,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}

	NewSuccessResponse(http.StatusOK, fmt.Sprintf("book with id %v updated sucessfully", updatedId), nil).RenderJSON(w)
}

func (b *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	if param == "" {
		NewErrorResponse("error to delete book",
			"the requested endpoint requires the query id parameter",
			http.StatusBadRequest,
			nil, r.URL.Path).
			RenderJSON(w)
		return
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		NewErrorResponse("error to delete book",
			"parameter with invalid format",
			http.StatusBadRequest,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}

	err = b.BookService.Delete(uint(id))
	if err != nil {
		NewErrorResponse("error to delete book",
			"",
			http.StatusBadRequest,
			err, r.URL.Path).
			RenderJSON(w)
	}
	NewSuccessResponse(http.StatusOK, fmt.Sprintf("book with id %v updated sucessfully", id), nil).RenderJSON(w)
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
		NewErrorResponse("error to list books",
			"error to parse queries",
			http.StatusBadRequest,
			err, r.URL.Path).
			RenderJSON(w)
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
		NewErrorResponse("error to list books",
			"",
			http.StatusBadRequest,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}

	result, err := b.BookService.List(*pag, aux)
	if err != nil {
		NewErrorResponse("error to list books",
			"",
			http.StatusBadRequest,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}
	NewSuccessResponse(http.StatusOK, "", result).RenderJSON(w)
}
