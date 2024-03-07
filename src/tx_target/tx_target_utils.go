package txtarget

import (
	"context"
	"math/big"

	"github.com/ledgerwatch/erigon-lib/common"
)

func (tt *TxTarget) getBalanceByHexAddress(address string) (*big.Int, error) {
	result, err := tt.ethClient.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		return nil, err
	}
	return result.ToBig(), nil
}

func (tt *TxTarget) shouldApplyTxLimit() bool {
	return tt.txSendingLimit > 0
}
