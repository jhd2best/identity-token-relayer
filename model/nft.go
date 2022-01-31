package model

type Nft struct {
	ContractAddress  string `firestore:"contract_address"`
	TokenId          int64  `firestore:"token_id"`
	OwnerAddress     string `firestore:"owner_address"`
	LastUpdateHeight int64  `firestore:"last_update_height"`
}
