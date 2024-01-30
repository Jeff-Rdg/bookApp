package handlers

import (
	"BookApp/internal/entities/book"
	"BookApp/internal/httpResponse"
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
		httpResponse.NewBadRequestError("error to find book by id", "no parameter informed", r).RenderJSON(w)
		return
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		httpResponse.NewBadRequestError("error to find book by id", err.Error(), r).RenderJSON(w)
		return
	}

	result, err := b.BookService.GetById(uint(id))
	if err != nil {
		httpResponse.NewNotFoundError("error to find book by id", err.Error(), r).RenderJSON(w)
		return
	}
	httpResponse.NewSuccessResponse(result).RenderJSON(w)
}

func (b *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newBook book.Request

	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		httpResponse.NewBadRequestError("error to create new book", err.Error(), r).RenderJSON(w)
		return
	}

	id, err := b.BookService.Create(newBook)
	if err != nil {
		httpResponse.NewBadRequestError("error to create new book", err.Error(), r).RenderJSON(w)
		return
	}

	httpResponse.NewCreatedResponse(fmt.Sprintf("book with id %v created successfully", id)).RenderJSON(w)
}

func (b *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	var bookRequest book.Request
	param := chi.URLParam(r, "id")
	if param == "" {
		httpResponse.NewBadRequestError("error to update book", "no parameter informed", r).RenderJSON(w)
		return
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		httpResponse.NewBadRequestError("error to update book", err.Error(), r).RenderJSON(w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&bookRequest)
	if err != nil {
		httpResponse.NewBadRequestError("error to update book", err.Error(), r).RenderJSON(w)
		return
	}

	updatedId, err := b.BookService.Update(uint(id), bookRequest)
	if err != nil {
		httpResponse.NewBadRequestError("error to update book", err.Error(), r).RenderJSON(w)
		return
	}

	httpResponse.NewSuccessResponse(fmt.Sprintf("book with id %v updated sucessfully", updatedId)).RenderJSON(w)
}

func (b *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	if param == "" {
		httpResponse.NewBadRequestError("error to delete book", "no parameter informed", r).RenderJSON(w)
		return
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		httpResponse.NewBadRequestError("error to delete book", err.Error(), r).RenderJSON(w)
		return
	}

	err = b.BookService.Delete(uint(id))
	if err != nil {
		httpResponse.NewInternalServerError("error to delete book", err.Error(), r).RenderJSON(w)
	}
	httpResponse.NewSuccessResponse(fmt.Sprintf("book with id %v deleted sucessfully", id)).RenderJSON(w)
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
		httpResponse.NewBadRequestError("error to list books", err.Error(), r).RenderJSON(w)
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
		httpResponse.NewBadRequestError("error to list books", err.Error(), r).RenderJSON(w)
		return
	}

	result, err := b.BookService.List(*pag, aux)
	if err != nil {
		httpResponse.NewBadRequestError("error to list books", err.Error(), r).RenderJSON(w)
		return
	}

	httpResponse.NewSuccessResponse(result).RenderJSON(w)
}
