package internal

import (
	"github.com/ostafen/clover/v2"
)

type pricingData struct {
	Value        float64
	CurrencyCode string
}

func getPricingData(id string, db *clover.DB) (*pricingData, error) {
	doc, err := db.FindFirst(clover.NewQuery(pricingCollection).Where(clover.Field("productId").Eq(id)))
	if err != nil {
		return nil, err
	}

	if doc == nil {
		return &pricingData{}, nil
	}

	priceData := &productPriceDocument{}
	err = doc.Unmarshal(priceData)
	if err != nil {
		return nil, err
	}

	return &pricingData{
		Value:        priceData.Value,
		CurrencyCode: priceData.CurrencyCode,
	}, nil
}

func updatePricingData(id string, price *Price, db *clover.DB) error {
	existingDoc, err := db.FindFirst(clover.NewQuery(pricingCollection).Where(clover.Field("productId").Eq(id)))
	if err != nil {
		return err
	}

	if existingDoc == nil {
		return insertNewPricingData(id, price, db)
	}

	updates := make(map[string]interface{})
	updates["productId"] = id
	updates["value"] = price.Value
	updates["currencyCode"] = price.CurrencyCode

	return db.UpdateById(pricingCollection, existingDoc.ObjectId(), updates)
}

func insertNewPricingData(id string, price *Price, db *clover.DB) error {
	newDoc := clover.NewDocument()
	newDoc.Set("productId", id)
	newDoc.Set("value", price.Value)
	newDoc.Set("currencyCode", price.CurrencyCode)

	return db.Insert(pricingCollection, newDoc)
}
