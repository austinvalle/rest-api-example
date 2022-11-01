package internal

import (
	"fmt"
	"log"
	"os"

	"github.com/ostafen/clover/v2"
)

const (
	pricingCollection = "productPrices"
	dbSeedPathEnv     = "DB_SEED_PATH"
)

type productPriceDocument struct {
	ProductId    string  `clover:"productId"`
	Value        float64 `clover:"value"`
	CurrencyCode string  `clover:"currencyCode"`
}

func InitDB() (*clover.DB, error) {
	db, err := clover.Open("db/pricing_database")
	if err != nil {
		return nil, err
	}

	exists, _ := db.Exists(clover.NewQuery(pricingCollection))

	if !exists {
		err = seedDB(db)
		if err != nil {
			log.Printf("error during seed: %s - creating empty collection...\n", err)
			collectionErr := db.CreateCollection(pricingCollection)
			if collectionErr != nil {
				return nil, fmt.Errorf("error initializing collection - %w", collectionErr)
			}
		}
	}

	return db, nil
}

func seedDB(db *clover.DB) error {
	seedPath := os.Getenv(dbSeedPathEnv)
	if seedPath == "" {
		seedPath = "db/pricing-seed.json"
		log.Printf("'%s' environment variable not set, attempting to seed with '%s'\n", dbSeedPathEnv, seedPath)
	}
	return db.ImportCollection(pricingCollection, seedPath)
}
