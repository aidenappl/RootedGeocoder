package structs

type Organisation struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	LocationID   int    `json:"location_id"`
	AddressLine1 string `json:"address_line_1"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZipCode      string `json:"zip_code"`
}
