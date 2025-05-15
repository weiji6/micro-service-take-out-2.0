package request

type RegisterRequest struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

type PayRequest struct {
	UserID   int     `json:"user_id"`
	ItemID   int     `json:"item_id"`
	Quantity int     `json:"quantity"`
	Amount   float32 `json:"amount"`
}

type ItemRequest struct {
	ItemID   int     `json:"item_id"`
	ItemName string  `json:"item_name"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
}
