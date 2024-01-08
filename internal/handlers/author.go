package handlers

import (
	"BookApp/author"
	"BookApp/pkg/filter"
	"BookApp/pkg/pagination"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

type AuthorHandler struct {
	AuthorService author.UseCase
}

func NewAuthorHandler(service author.UseCase) *AuthorHandler {
	return &AuthorHandler{AuthorService: service}
}

func (a *AuthorHandler) FindById(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	if param == "" {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "No id informed"})
		return
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "error to find author", Data: err})
		return
	}

	result, err := a.AuthorService.GetById(uint(id))
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "error to find author", Data: err})
		return
	}

	RenderJSON(w, Result{StatusCode: http.StatusOK, Data: result})
}

func (a *AuthorHandler) UploadCsv(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "error to analisys form", Data: err})
		return
	}

	file, _, err := r.FormFile("csv_file")
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "error to get CSV", Data: err})
		return
	}

	defer file.Close()

	err = a.AuthorService.ReadCsv(file, func(request author.Request) error {
		_, err = a.AuthorService.Create(request)
		return err
	})

	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "error to process CSV", Data: err})
		return
	}

	RenderJSON(w, Result{StatusCode: http.StatusOK, Message: "csv uploaded successfully"})
}

func (a *AuthorHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	sort := r.URL.Query().Get("sort")
	param := r.URL.Query().Get("name")

	aux := filter.Author{Name: param}

	pag, err := pagination.NewPagination(limit, page, sort)
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "error to parse queries", Data: err})
		return
	}

	result, err := a.AuthorService.List(*pag, aux)
	if err != nil {
		RenderJSON(w, Result{StatusCode: http.StatusBadRequest, Message: "error to list authors", Data: err})
		return
	}

	RenderJSON(w, Result{StatusCode: http.StatusOK, Data: result})
}
