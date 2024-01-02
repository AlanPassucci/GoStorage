package domain

type Sale struct {
	Id         int `json:"id"`
	ProductId  int `json:"product_id"`
	InvoicesId int `json:"invoice_id"`
	Quantity   int `json:"quantity"`
}
