package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/austinvalle/rest-api-example/internal"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// TODO: handle error
	db, err := internal.InitDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close()

	productEndpoint := productEndpoint{db: db}

	r.Get("/products/{id:\\d+}", productEndpoint.getProductById)
	r.Put("/products/{id:\\d+}/price", productEndpoint.updateProductPriceById)

	// TODO: add env variable
	port := 3000
	fmt.Printf("API is running at http://localhost:%d\n", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
