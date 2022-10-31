package internal

import (
	"strconv"

	"github.com/ostafen/clover/v2"
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

func GetProduct(id string, db *clover.DB) (*Product, error) {
	productData, err := getProductDataFromExternal(id)
	if err != nil {
		return nil, err
	}

	pricingData, err := getPricingData(id, db)
	if err != nil {
		return nil, err
	}

	return mergeProductAndPriceData(id, productData, pricingData), nil
}

func UpdateProduct(id string, price *Price, db *clover.DB) error {
	// Check if product exists before updating price
	_, err := getProductDataFromExternal(id)
	if err != nil {
		return err
	}

	return updatePricingData(id, price, db)
}

func mergeProductAndPriceData(id string, productData *productData, pricingData *pricingData) *Product {
	return &Product{
		Id:   id,
		Name: productData.Item.ProductDescription.Title,
		CurrentPrice: Price{
			Value:        Currency(pricingData.Value),
			CurrencyCode: pricingData.CurrencyCode,
		},
	}
}
