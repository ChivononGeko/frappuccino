package models

type TotalPrice struct {
	TotalSale float64 `json:"total-sales"`
}

type PopularItem struct {
	ItemName        string `json:"item-name"`
	QuantityOfSales int    `json:"quantity_of_sales"`
}

func NewTotalPrice() *TotalPrice {
	return &TotalPrice{}
}

func NewPopularItem(name string, quantity int) PopularItem {
	return PopularItem{
		ItemName:        name,
		QuantityOfSales: quantity,
	}
}
