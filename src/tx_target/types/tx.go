package types

import (
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	txsourcetypes "github.com/gateway-fm/tx-repeater/src/tx_source/types"
	"github.com/gateway-fm/tx-repeater/src/utils"
)

type Tx struct {
	Bytes  []byte  `json:"bytes"`
	From   string  `json:"from"`
	To     *string `json:"to,omitempty"`
	Hash   string  `json:"hash"`
	Status uint64  `json:"status"`
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
		Bytes:  bytes,
		From:   from,
		To:     to,
		Hash:   hash,
		Status: status,
	}
}

func FromSourceTx(resp *txsourcetypes.Tx, txReceipt *txsourcetypes.TxReceipt) (*Tx, error) {
	t := resp.Result

	gasPrice := utils.StringToBig(t.GasPrice)
	value := utils.StringToBig(t.Value)
	v := utils.StringToBig(t.V)
	r := utils.StringToBig(t.R)
	s := utils.StringToBig(t.S)

	gas, err := utils.StringToUint(t.Gas)
	if err != nil {
		return nil, err
	}

	nonce, err := utils.StringToUint(t.Nonce)
	if err != nil {
		return nil, err
	}

	var to []byte
	if t.To != nil {
		to, err = hex.DecodeString(utils.TrimHex(*t.To))
		if err != nil {
			return nil, err
		}
	}

	input, err := hex.DecodeString(utils.TrimHex(t.Input))
	if err != nil {
		return nil, err
	}

	toAddr := common.HexToAddress(hex.EncodeToString(to))
	toAddrPointer := &toAddr
	if t.To == nil {
		toAddrPointer = nil
	}
	txData := &ethtypes.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gas,
		To:       toAddrPointer,
		Value:    value,
		Data:     input,
		V:        v,
		R:        r,
		S:        s,
	}
	tx := ethtypes.NewTx(txData)
	b, _ := rlp.EncodeToBytes(tx)
	return NewTx(b, t.From, t.To, t.Hash, txReceipt.ParseAndGetStatus()), nil
}

func (tx *Tx) IsToBridgeTx() bool {
	if tx.To == nil {
		return false
	}

	return strings.EqualFold(*tx.To, utils.BRIDGE_ADDRESS)
}
