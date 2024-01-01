package main

import (
	"BookApp/configs"
	"BookApp/internal/routes"
	"net/http"
)

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := configs.LoadDatabase()
	if err != nil {
		panic(err)
	}

	r := routes.LoadRoutes(db)

	http.ListenAndServe(":8000", r)
}
