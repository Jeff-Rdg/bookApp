package routes

import (
	"BookApp/internal/entities/author"
	"BookApp/internal/entities/book"
	"BookApp/internal/handlers"
	"BookApp/internal/middlewares"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/gorm"
)

func LoadRoutes(db *gorm.DB) *chi.Mux {
	authorService := author.NewService(db)
	authorHandler := handlers.NewAuthorHandler(authorService)

	bookService := book.NewService(db)
	bookHandler := handlers.NewBookHandler(bookService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.LogException)

	r.Route("/author", func(r chi.Router) {
		r.Get("/", authorHandler.List)
		r.Post("/upload_csv", authorHandler.UploadCsv)
		r.Get("/{id:[0-9]+}", authorHandler.FindById)
		r.Get("/search/{name}", authorHandler.FindByName)
	})

	r.Route("/book", func(r chi.Router) {
		r.Get("/", bookHandler.List)
		r.Get("/{id}", bookHandler.FindById)
		r.Post("/", bookHandler.Create)
		r.Put("/{id}", bookHandler.Update)
		r.Delete("/{id}", bookHandler.Delete)
	})

	return r
}
