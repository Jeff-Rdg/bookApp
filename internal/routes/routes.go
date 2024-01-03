package routes

import (
	"BookApp/author"
	"BookApp/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/gorm"
)

func LoadRoutes(db *gorm.DB) *chi.Mux {
	authorService := author.NewService(db)
	authorHandler := handlers.NewAuthorHandler(authorService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/author", func(r chi.Router) {
		r.Get("/{id}", authorHandler.FindById)
		r.Get("/list", authorHandler.List)
		r.Post("/upload_csv", authorHandler.UploadCsv)
	})

	return r
}
