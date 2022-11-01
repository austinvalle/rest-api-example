package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/austinvalle/rest-api-example/internal"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const (
	apiPortEnv  = "API_PORT"
	defaultPort = 3000
)

func init() {
	_, exists := os.LookupEnv(internal.ApiKeyEnv)
	if !exists {
		log.Fatalf("missing environment variable '%s' - exiting...", internal.ApiKeyEnv)
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	db, err := internal.InitDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close()

	productEndpoint := productEndpoint{db: db}

	r.Get("/products/{id:\\d+}", productEndpoint.getProductById)
	r.Put("/products/{id:\\d+}/price", productEndpoint.updateProductPriceById)

	port := getPort()
	log.Printf("API is running at http://localhost:%d\n", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func getPort() int {
	portStr := os.Getenv(apiPortEnv)
	if portStr == "" {
		return defaultPort
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Invalid '%s' environment variable: %s - defaulting to '%d'", apiPortEnv, portStr, defaultPort)
		return defaultPort
	}

	return port
}
