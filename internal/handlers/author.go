package handlers

import (
	"BookApp/internal/entities/author"
	"BookApp/internal/httpResponse"
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
		httpResponse.NewBadRequestError("error to find author by id", "no parameter informed", r).RenderJSON(w)
		return
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		httpResponse.NewBadRequestError("error to find author by id", err.Error(), r).RenderJSON(w)
		return
	}

	result, err := a.AuthorService.GetById(uint(id))
	if err != nil {
		httpResponse.NewNotFoundError("error to find author by id", err.Error(), r).RenderJSON(w)
		return
	}

	httpResponse.NewSuccessResponse(result).RenderJSON(w)
}

func (a *AuthorHandler) FindByName(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "name")
	if param == "" {
		httpResponse.NewBadRequestError("error to find author by name", "no parameter informed", r).RenderJSON(w)
		return
	}

	result, err := a.AuthorService.GetByName(param)
	if err != nil {
		httpResponse.NewNotFoundError("error to find author by name", err.Error(), r).RenderJSON(w)
		return
	}

	httpResponse.NewSuccessResponse(result).RenderJSON(w)
}

func (a *AuthorHandler) UploadCsv(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		httpResponse.NewBadRequestError("error to Upload csv file", err.Error(), r).RenderJSON(w)
		return
	}

	file, _, err := r.FormFile("csv_file")
	if err != nil {
		httpResponse.NewBadRequestError("error to read csv file", err.Error(), r).RenderJSON(w)
		return
	}

	defer file.Close()

	err = a.AuthorService.ReadCsv(file, func(request author.Request) error {
		aut, _ := author.NewAuthor(request)
		_, err = a.AuthorService.Create(aut)
		return err
	})

	if err != nil {
		httpResponse.NewBadRequestError("error to process csv file", err.Error(), r).RenderJSON(w)
		return
	}

	httpResponse.NewCreatedResponse("csv uploaded successfully").RenderJSON(w)
}

func (a *AuthorHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	sort := r.URL.Query().Get("sort")
	param := r.URL.Query().Get("name")

	aux := filter.Author{Name: param}

	pag, err := pagination.NewPagination(limit, page, sort)
	if err != nil {
		httpResponse.NewBadRequestError("error to list authors", err.Error(), r).RenderJSON(w)
		return
	}

	result, err := a.AuthorService.List(*pag, aux)
	if err != nil {
		httpResponse.NewBadRequestError("error to list authors", err.Error(), r).RenderJSON(w)
		return
	}

	httpResponse.NewSuccessResponse(result).RenderJSON(w)
}
