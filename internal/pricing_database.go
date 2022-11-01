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

type ProductPriceDocument struct {
	ProductId    string  `clover:"productId"`
	Value        float64 `clover:"value"`
	CurrencyCode string  `clover:"currencyCode"`
}

type PricingDatabase interface {
	GetPricingById(id string) (*ProductPriceDocument, error)
	UpsertPricing(ProductPriceDocument) error
	InsertPricing(ProductPriceDocument) error
	Close() error
}

type pricingDatabase struct {
	db *clover.DB
}

func (pdb pricingDatabase) GetPricingById(id string) (*ProductPriceDocument, error) {
	db := pdb.db

	doc, err := db.FindFirst(clover.NewQuery(pricingCollection).Where(clover.Field("productId").Eq(id)))
	if err != nil {
		return nil, err
	}

	// TODO: feels weird
	if doc == nil {
		return &ProductPriceDocument{}, nil
	}

	priceData := &ProductPriceDocument{}
	err = doc.Unmarshal(priceData)
	if err != nil {
		return nil, err
	}

	return priceData, nil
}

func (pdb pricingDatabase) UpsertPricing(doc ProductPriceDocument) error {
	db := pdb.db

	existingDoc, err := db.FindFirst(clover.NewQuery(pricingCollection).Where(clover.Field("productId").Eq(doc.ProductId)))
	if err != nil {
		return err
	}

	if existingDoc == nil {
		return pdb.InsertPricing(doc)
	}

	updates := make(map[string]interface{})
	updates["productId"] = doc.ProductId
	updates["value"] = doc.Value
	updates["currencyCode"] = doc.CurrencyCode

	return db.UpdateById(pricingCollection, existingDoc.ObjectId(), updates)
}

func (pdb pricingDatabase) InsertPricing(doc ProductPriceDocument) error {
	db := pdb.db

	cloverDoc := clover.NewDocument()
	cloverDoc.Set("productId", doc.ProductId)
	cloverDoc.Set("value", doc.Value)
	cloverDoc.Set("currencyCode", doc.CurrencyCode)

	return db.Insert(pricingCollection, cloverDoc)
}

func (pdb pricingDatabase) Close() error {
	return pdb.db.Close()
}

func InitDB() (PricingDatabase, error) {
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

	database := &pricingDatabase{db: db}

	return database, nil
}

func seedDB(db *clover.DB) error {
	seedPath := os.Getenv(dbSeedPathEnv)
	if seedPath == "" {
		seedPath = "db/pricing-seed.json"
		log.Printf("'%s' environment variable not set, attempting to seed with '%s'\n", dbSeedPathEnv, seedPath)
	}
	return db.ImportCollection(pricingCollection, seedPath)
}
