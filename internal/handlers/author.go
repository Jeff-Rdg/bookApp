package handlers

import (
	"BookApp/author"
	"BookApp/pkg"
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
		RenderJSON(w, http.StatusBadRequest, "no id informed", nil)
		return
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		RenderJSON(w, http.StatusBadRequest, "error to find author", err)
		return
	}

	result, err := a.AuthorService.GetById(uint(id))
	if err != nil {
		RenderJSON(w, http.StatusBadRequest, "error to find author", err)
		return
	}

	RenderJSON(w, http.StatusOK, "", result)
}

func (a *AuthorHandler) UploadCsv(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		RenderJSON(w, http.StatusBadRequest, "error to analisys form", err)
		return
	}

	file, _, err := r.FormFile("csv_file")
	if err != nil {
		RenderJSON(w, http.StatusBadRequest, "error to get CSV", err)
		return
	}

	defer file.Close()

	err = a.AuthorService.ReadCsv(file, func(request author.Request) error {
		_, err = a.AuthorService.Create(request)
		return err
	})

	if err != nil {
		RenderJSON(w, http.StatusBadRequest, "error to process CSV", err)
		return
	}

	RenderJSON(w, http.StatusOK, "csv uploaded successfully", nil)
}

func (a *AuthorHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	sort := r.URL.Query().Get("sort")

	pagination, err := pkg.NewPagination(limit, page, sort)
	if err != nil {
		RenderJSON(w, http.StatusBadRequest, "error to parse queries", err)
		return
	}

	result, err := a.AuthorService.List(*pagination)
	if err != nil {
		RenderJSON(w, http.StatusBadRequest, "error to list authors", err)
		return
	}

	RenderJSON(w, http.StatusOK, "", result)
}
