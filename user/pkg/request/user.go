package request

type UserRequest struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
}
