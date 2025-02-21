package models

// Order adalah struktur data untuk pesanan
type Order struct {
	Product       string `json:"product"`
	Quantity      int    `json:"quantity"`
	UnitPrice     int    `json:"unit_price"`
	ShippingCost  int    `json:"shipping_cost"`
	InsuranceCost int    `json:"insurance_cost"`
	TotalPrice    int    `json:"total_price"`
	FullName      string `json:"full_name"`
	Address       string `json:"address"`
	PostalCode    string `json:"postal_code"`
	Phone         string `json:"phone"`
}
