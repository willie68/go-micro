package pmodel

// Address this is the main model
type Address struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Firstname string `json:"firstname"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zip_code"`
}
