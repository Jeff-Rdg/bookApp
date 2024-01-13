package routes

import (
	"BookApp/author"
	"BookApp/book"
	"BookApp/internal/handlers"
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

	r.Route("/author", func(r chi.Router) {
		r.Get("/{id}", authorHandler.FindById)
		r.Get("/{name}", authorHandler.FindByName)
		r.Get("/list", authorHandler.List)
		r.Post("/upload_csv", authorHandler.UploadCsv)
	})

	r.Route("/book", func(r chi.Router) {
		r.Get("/{id}", bookHandler.FindById)
		r.Get("/list", bookHandler.List)
		r.Post("/", bookHandler.Create)
		r.Put("/{id}", bookHandler.Update)
	})

	return r
}
