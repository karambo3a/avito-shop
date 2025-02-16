package requests

type SendCoinRequest struct {
	ToUser string `json:"toUser" binding:"required"`
	Amount int `json:"amount" binding:"required"`
}
