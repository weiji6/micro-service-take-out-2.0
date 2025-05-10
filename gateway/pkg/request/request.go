package request

type RegisterRequest struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

type PayRequest struct {
	UserID int     `json:"user_id"`
	ItemID int     `json:"item_id"`
	Amount float32 `json:"amount"`
}
