package types

type Block struct {
	Number       int   `json:"number"`
	Transactions []*Tx `json:"transactions"`
}

func NewBlock(number, size int) *Block {
	return &Block{
		Number:       number,
		Transactions: make([]*Tx, 0, size),
	}
}

func (b *Block) HasTransactions() bool {
	return len(b.Transactions) > 0
}

func (b *Block) AppendTx(tx *Tx) {
	b.Transactions = append(b.Transactions, tx)
}
