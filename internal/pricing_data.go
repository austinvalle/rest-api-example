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
