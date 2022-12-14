package internal

import (
	"strconv"
)

type Product struct {
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	CurrentPrice Price  `json:"current_price,omitempty"`
}

type Currency float64
type Price struct {
	Value        Currency `json:"value,omitempty"`
	CurrencyCode string   `json:"currency_code,omitempty"`
}

// Outputs float JSON with 2 digit precision
func (c Currency) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(c), 'f', 2, 64)), nil
}

func GetProduct(id string, db PricingDatabase) (*Product, error) {
	productData, err := getProductDataFromExternal(id)
	if err != nil {
		return nil, err
	}

	pricingData, err := db.GetPricingById(id)
	if err != nil {
		return nil, err
	}

	return mergeProductAndPriceData(id, productData, pricingData), nil
}

func UpdateProduct(id string, price *Price, db PricingDatabase) error {
	// Check if product exists before updating price
	_, err := getProductDataFromExternal(id)
	if err != nil {
		return err
	}

	return db.UpsertPricing(ProductPriceDocument{
		ProductId:    id,
		Value:        float64(price.Value),
		CurrencyCode: price.CurrencyCode,
	})
}

func mergeProductAndPriceData(id string, productData *productData, pricingData *ProductPriceDocument) *Product {
	return &Product{
		Id:   id,
		Name: productData.Item.ProductDescription.Title,
		CurrentPrice: Price{
			Value:        Currency(pricingData.Value),
			CurrencyCode: pricingData.CurrencyCode,
		},
	}
}
