package domain

type Carries struct {
	ID int `json:"id"`
	CID string `json:"cid"`
	CompanyName string `json:"company_name"`
	Address string `json:"address"`
	Telephone string `json:"telephone"`
	LocalityID int `json:"locality_id"`
}

type LocalityCarries struct {
	LocalityID string `json:"locality_id"`
	LocalityName string `json:"locality_name"`
	CarriesCount int `json:"carries_count"`
}