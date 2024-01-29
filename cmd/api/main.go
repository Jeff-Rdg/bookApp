package main

import (
	"BookApp/configs"
	"BookApp/internal/routes"
	"fmt"
	"net/http"
)

func main() {
	conf, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := configs.LoadDatabase(conf)
	if err != nil {
		panic(err)
	}

	r := routes.LoadRoutes(db)
	port := fmt.Sprintf(":%s", conf.WebServerPort)
	if port == ":" {
		http.ListenAndServe(":8000", r)
	} else {
		http.ListenAndServe(port, r)
	}
}
