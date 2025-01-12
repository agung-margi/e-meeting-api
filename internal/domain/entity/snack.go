package entity

type Snack struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Category  string `json:"category"`
	Price     int    `json:"price"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
