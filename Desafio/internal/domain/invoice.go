package domain

type Invoice struct {
	Id         int     `json:"id"`
	Datetime   string  `json:"datetime"`
	CustomerId int     `json:"customer_id"`
	Total      float64 `json:"total"`
}
