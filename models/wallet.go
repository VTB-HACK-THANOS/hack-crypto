package models

type WalletCredentials struct {
	// {
	// 	"privateKey": "f65eb3320a3e9dd51eaf2d67d71e609459081cf63aefd32fda79c425c296f257",
	// 	"publicKey": "0x15Cc4abzz27647ec9fE70D892E55586074263dF0"
	// }
	PrivateKey string `json:"privateKey" bun:"private_key"`
	PublicKey  string `json:"publicKey" bun:"public_key"`
}
