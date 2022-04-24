package tools

type NftOwnerResponse struct {
	Total  int         `json:"total"`
	Page   int         `json:"page"`
	Status string      `json:"status"`
	Cursor string      `json:"cursor"`
	Result []OwnerItem `json:"result"`
}

type OwnerItem struct {
	TokenId     string `json:"token_id"`
	OwnerOf     string `json:"owner_of"`
	BlockNumber string `json:"block_number"`
}
