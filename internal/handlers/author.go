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
		MakeErrorResponse("error to find author by id",
			"the requested endpoint requires the query id parameter",
			http.StatusBadRequest,
			nil, r.URL.Path).
			RenderJSON(w)
		return
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		MakeErrorResponse("error to find author by id",
			"parameter with invalid format",
			http.StatusBadRequest,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}

	result, err := a.AuthorService.GetById(uint(id))
	if err != nil {
		MakeErrorResponse("error to find author by id",
			"",
			http.StatusNotFound,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}
	MakeSuccessResponse(http.StatusOK, "", result).RenderJSON(w)
}

func (a *AuthorHandler) FindByName(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "name")
	if param == "" {
		MakeErrorResponse("error to find author by name",
			"the requested endpoint requires the query name parameter",
			http.StatusBadRequest,
			nil, r.URL.Path).
			RenderJSON(w)
		return
	}

	result, err := a.AuthorService.GetByName(param)
	if err != nil {
		MakeErrorResponse("error to find author by name",
			"",
			http.StatusNotFound,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}

	MakeSuccessResponse(http.StatusOK, "", result).RenderJSON(w)
}

func (a *AuthorHandler) UploadCsv(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		MakeErrorResponse("error to Upload csv file",
			"",
			http.StatusBadRequest,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}

	file, _, err := r.FormFile("csv_file")
	if err != nil {
		MakeErrorResponse("error to read csv file",
			"",
			http.StatusBadRequest,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}

	defer file.Close()

	err = a.AuthorService.ReadCsv(file, func(request author.Request) error {
		_, err = a.AuthorService.Create(request)
		return err
	})

	if err != nil {
		MakeErrorResponse("error to process csv file",
			"",
			http.StatusBadRequest,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}

	MakeSuccessResponse(http.StatusOK, "csv uploaded successfully", nil).RenderJSON(w)
}

func (a *AuthorHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	sort := r.URL.Query().Get("sort")
	param := r.URL.Query().Get("name")

	aux := filter.Author{Name: param}

	pag, err := pagination.NewPagination(limit, page, sort)
	if err != nil {
		MakeErrorResponse("error to list authors",
			"",
			http.StatusBadRequest,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}

	result, err := a.AuthorService.List(*pag, aux)
	if err != nil {
		MakeErrorResponse("error to list authors",
			"",
			http.StatusBadRequest,
			err, r.URL.Path).
			RenderJSON(w)
		return
	}

	MakeSuccessResponse(http.StatusOK, "", result).RenderJSON(w)
}
