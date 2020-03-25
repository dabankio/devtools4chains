package omnitypes

type SenddissuancefixedCmd struct {
	Fromaddress                                    string
	Ecosystem                                      int
	Typ                                            int
	Previousid                                     int
	Category, Subcategory, Name, URL, Data, Amount string
}

// GettransactionResult .
type GettransactionResult struct {
	Txid             string `json:"txid"`             // (string) the hex-encoded hash of the transaction
	Sendingaddress   string `json:"sendingaddress"`   // (string) the Bitcoin address of the sender
	Referenceaddress string `json:"referenceaddress"` // (string) a Bitcoin address used as reference (if any)
	Ismine           bool   `json:"ismine"`           // (boolean) whether the transaction involes an address in the wallet
	Confirmations    int    `json:"confirmations"`    // (number) the number of transaction confirmations
	Fee              string `json:"fee"`              // (string) the transaction fee in bitcoins
	Blocktime        int    `json:"blocktime"`        // (number) the timestamp of the block that contains the transaction
	Valid            bool   `json:"valid"`            // (boolean) whether the transaction is valid
	Positioninblock  int    `json:"positioninblock"`  // (number) the position (index) of the transaction within the block
	Version          int    `json:"version"`          // (number) the transaction version
	TypeInt          int    `json:"type_int"`         // (number) the transaction type as number
	Type             string `json:"type"`             // (string) the transaction type as string
	//other
	Propertyid int `json:"propertyid"`
}

// GetbalanceResult .
type GetbalanceResult struct {
	Balance  string `json:"balance"`
	Reserved string `json:"reserved"`
	Frozen   string `json:"frozen"`
}
