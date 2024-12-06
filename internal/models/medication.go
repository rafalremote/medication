package models

type Medication struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Dosage    string `json:"dosage"`
	Form      string `json:"form"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
