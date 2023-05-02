package example


type Input struct {
	TxID    string
	Index   uint32
	Address string
	Value   uint64
	Wif     string
}

type Output struct {
	Address string
	Value   uint64
}

type Utxo struct {
	Address string
	Value   uint64
	Index   uint64
	Tx      string
}