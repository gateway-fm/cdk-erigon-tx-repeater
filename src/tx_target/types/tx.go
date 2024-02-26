package types

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	txsourcetypes "github.com/gateway-fm/tx-repeater/src/tx_source/types"
)

type Tx struct {
	Bytes []byte `json:"bytes"`
	From  string `json:"from"`
	Hash  string `json:"hash"`
}

type TxType int64

const (
	LegacyTx     TxType = 0x0
	AccessListTx TxType = 0x01
	StateTx      TxType = 0x7f
	DynamicFeeTx TxType = 0x02
)

func NewTx(bytes []byte, from, hash string) *Tx {
	return &Tx{
		Bytes: bytes,
		From:  from,
		Hash:  hash,
	}
}

func FromSourceTx(resp *txsourcetypes.Tx) (*Tx, error) {
	t := resp.Result

	gasPrice := stringToBig(t.GasPrice)
	value := stringToBig(t.Value)
	v := stringToBig(t.V)
	r := stringToBig(t.R)
	s := stringToBig(t.S)

	gas, err := stringToUint(t.Gas)
	if err != nil {
		return nil, err
	}

	nonce, err := stringToUint(t.Nonce)
	if err != nil {
		return nil, err
	}

	var to []byte
	if t.To != nil {
		to, err = hex.DecodeString(trimHex(*t.To))
		if err != nil {
			return nil, err
		}
	}

	input, err := hex.DecodeString(trimHex(t.Input))
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
	return NewTx(b, t.From, t.Hash), nil
}

func stringToInt(input string) (int64, error) {
	return strconv.ParseInt(trimHex(input), 16, 64)
}

func stringToUint(input string) (uint64, error) {
	return strconv.ParseUint(trimHex(input), 16, 64)
}

func stringToBig(input string) *big.Int {
	result := new(big.Int)
	result.SetString(trimHex(input), 16)
	return result
}

func trimHex(input string) string {
	if strings.HasPrefix(input, "0x") {
		return input[2:]
	}
	return input
}
