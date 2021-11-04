package omnijson

type SignRawTransactionWithWalletResult = struct {
	Hex      string                    `json:"hex"`
	Complete bool                      `json:"complete"`
}

type SignRawTransactionWithWalletCommand struct {
	Hex      string
}

func (SignRawTransactionWithWalletCommand) Method() string {
	return "signrawtransactionwithwallet"
}

func (SignRawTransactionWithWalletCommand) ID() string {
	return "1"
}

func (cmd SignRawTransactionWithWalletCommand) Params() []interface{} {
	return []interface{}{cmd.Hex}
}
