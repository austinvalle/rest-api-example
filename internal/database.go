package internal

import (
	"fmt"

	"github.com/ostafen/clover/v2"
)

const (
	pricingCollection = "productPrices"
)

type productPriceDocument struct {
	ProductId    string  `clover:"productId"`
	Value        float64 `clover:"value"`
	CurrencyCode string  `clover:"currencyCode"`
}

func InitDB() (*clover.DB, error) {
	// TODO: Put location of DB in env variable?
	db, err := clover.Open("db/pricing_database")
	if err != nil {
		return nil, err
	}

	exists, _ := db.Exists(clover.NewQuery(pricingCollection))

	if !exists {
		err = seedDB(db)
		if err != nil {
			return nil, fmt.Errorf("failed to seed DB: %w", err)
		}
	}

	return db, nil
}

func seedDB(db *clover.DB) error {
	// TODO: Put location of seed in env variable?
	return db.ImportCollection(pricingCollection, "db/pricing-seed.json")
}
