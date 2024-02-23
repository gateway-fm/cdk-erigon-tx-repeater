package types

import (
	"strings"

	"github.com/gateway-fm/tx-repeater/src/utils"
)

type Tx struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		Nonce            string          `json:"nonce"`
		GasPrice         string          `json:"gasPrice"`
		GasTipCap        string          `json:"gasTipCap"`
		GasFeeCap        string          `json:"gasFeeCap"`
		Gas              string          `json:"gas"`
		To               *string         `json:"to,omitempty"`
		Value            string          `json:"value"`
		Input            string          `json:"input"`
		V                string          `json:"v"`
		R                string          `json:"r"`
		S                string          `json:"s"`
		Hash             string          `json:"hash"`
		From             string          `json:"from"`
		BlockHash        string          `json:"blockHash"`
		BlockNumber      string          `json:"blockNumber"`
		TransactionIndex string          `json:"transactionIndex"`
		Type             string          `json:"type"`
		ChainID          string          `json:"chainId"`
		AccessList       []TxAccessTuple `json:"accessList"`
	}
}

func (tx *Tx) IsBridgeTx() bool {
	return strings.ToLower(tx.Result.From) == strings.ToLower(utils.BRIDGE_ADDRESS)
}

type TxAccessTuple struct {
	Address     string
	StorageKeys []string
}
