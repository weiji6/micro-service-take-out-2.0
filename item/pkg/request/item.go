package request

type ItemRequest struct {
	ItemID   int     `json:"item_id"`
	ItemName string  `json:"item_name"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
}
