package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/austinvalle/rest-api-example/internal"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/ostafen/clover/v2"
)

type productEndpoint struct {
	db *clover.DB
}

type priceRequest struct {
	Value        *float64 `json:"value"`
	CurrencyCode *string  `json:"currency_code"`
}

func (req *priceRequest) Bind(r *http.Request) error {
	if req.Value == nil {
		return errors.New("missing required 'value' field")
	}
	if req.CurrencyCode == nil {
		return errors.New("missing required 'currency_code' field")
	}
	return nil
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

func (p *productEndpoint) updateProductPriceById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	priceReq, err := getPriceFromBody(r)
	if err != nil {
		render.Render(w, r, ErrBadRequest(fmt.Errorf("validation error - %w", err)))
		return
	}

	err = internal.UpdateProduct(id, priceReq, p.db)
	if errors.Is(err, internal.ErrProductNotFound) {
		render.Render(w, r, ErrNotFound(err))
		return
	} else if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.NoContent(w, r)
}

func getPriceFromBody(r *http.Request) (*internal.Price, error) {
	var priceReq priceRequest
	err := render.Bind(r, &priceReq)
	if err != nil {
		return nil, err
	}

	return &internal.Price{
		Value:        internal.Currency(*priceReq.Value),
		CurrencyCode: *priceReq.CurrencyCode,
	}, nil
}
