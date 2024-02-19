package domain

type Locality struct {
	ID           int    `json:"id"`
	PostalCode   int    `json:"postal_code"`
	LocalityName string `json:"locality_name"`
	ProvinceName string `json:"province_name"`
	CountryName  string `json:"country_name"`
}

type ReportSellers struct {
	Locality_id   int    `json:"locality_id"`
	Locality_name string `json:"locality_name"`
	Postal_code   int    `json:"postal_code"`
	Sellers_count int    `json:"sellers_count"`
}
