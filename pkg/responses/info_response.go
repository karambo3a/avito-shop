package responses

type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Recieved []Transaction `json:"received"`
	Sent     []Transaction `json:"sent"`
}

type Transaction interface{}

type SentTransaction struct {
	ToUser string `json:"to_user_id"`
	Amount int    `json:"amount"`
}

type RecievedTransaction struct {
	FromUser string `json:"from_user_id"`
	Amount   int    `json:"amount"`
}
