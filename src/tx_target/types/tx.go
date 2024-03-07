package types

import (
	"strings"

	"github.com/gateway-fm/tx-repeater/src/utils"
)

type Tx struct {
	Bytes []byte  `json:"bytes"`
	From  string  `json:"from"`
	To    *string `json:"to,omitempty"`
	Hash  string  `json:"hash"`
}

type TxType int64

const (
	LegacyTx     TxType = 0x0
	AccessListTx TxType = 0x01
	StateTx      TxType = 0x7f
	DynamicFeeTx TxType = 0x02
)

func NewTx(bytes []byte, from string, to *string, hash string, status uint64) *Tx {
	return &Tx{
		Bytes: bytes,
		From:  from,
		To:    to,
		Hash:  hash,
	}
}

func (tx *Tx) IsToBridgeTx() bool {
	if tx.To == nil {
		return false
	}

	return strings.EqualFold(*tx.To, utils.BRIDGE_ADDRESS)
}
