package main

import (
	"errors"
	"net/http"

	"github.com/austinvalle/rest-api-example/internal"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/ostafen/clover/v2"
)

type productEndpoint struct {
	db *clover.DB
}

func (p *productEndpoint) getProductById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	product, err := internal.GetProduct(id, p.db)
	if errors.Is(err, internal.ErrProductNotFound) {
		render.Render(w, r, ErrNotFound(err))
		return
	} else if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.JSON(w, r, product)
}
