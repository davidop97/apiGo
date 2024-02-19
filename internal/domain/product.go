package domain

// Product represents an underlying URL with statistics on how it is used.
type Product struct {
	ID             int     `json:"id"`
	Description    string  `json:"description"`
	ExpirationRate float32 `json:"expiration_rate"`
	FreezingRate   float32 `json:"freezing_rate"`
	Height         float32 `json:"height"`
	Length         float32 `json:"length"`
	Netweight      float32 `json:"netweight"`
	ProductCode    string  `json:"product_code"`
	RecomFreezTemp float32 `json:"recommended_freezing_temperature"`
	Width          float32 `json:"width"`
	ProductTypeID  int     `json:"product_type_id"`
	SellerID       int     `json:"seller_id"`
}

// Struct for the product record
type ProductRecord struct {
	ID            int     `json:"id"`
	LastUpdate    string  `json:"last_update_date"`
	PurchasePrice float32 `json:"purchase_price"`
	SalePrice     float32 `json:"sale_price"`
	ProductID     int     `json:"product_id"` //Product_code of the product (fk)
}

type ProductRecordCreate struct {
	LastUpdate    string  `json:"last_update_date"`
	PurchasePrice float32 `json:"purchase_price"`
	SalePrice     float32 `json:"sale_price"`
	ProductID     int     `json:"product_id"` //Product_code of the product (fk)
}

type ProductRecordGet struct {
	ProductID   int    `json:"product_id"` //Product_code of the product (fk)
	Description string `json:"description"`
	RecordCount int    `json:"record_count"` // count of records
}
