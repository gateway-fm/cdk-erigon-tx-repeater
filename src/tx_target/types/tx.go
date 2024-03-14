package types

import "github.com/holiman/uint256"

type TxType int64

const (
	LegacyTx     TxType = 0x0
	AccessListTx TxType = 0x01
	StateTx      TxType = 0x7f
	DynamicFeeTx TxType = 0x02
)

type Tx struct {
	Bytes   []byte  `json:"bytes"`
	From    string  `json:"from"`
	To      *string `json:"to,omitempty"`
	Hash    string  `json:"hash"`
	Funding string  `json:"funding"`
}

func NewTx(bytes []byte, from string, to *string, hash string, funding string) *Tx {
	return &Tx{
		Bytes:   bytes,
		From:    from,
		To:      to,
		Hash:    hash,
		Funding: funding,
	}
}

func (tx *Tx) GetFundingAsUint256() *uint256.Int {
	return uint256.MustFromHex(tx.Funding)
}
