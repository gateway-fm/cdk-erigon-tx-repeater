package txtarget

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	types "github.com/gateway-fm/tx-repeater/src/tx_target/types"
	"github.com/gateway-fm/tx-repeater/src/utils"
	"github.com/ledgerwatch/erigon-lib/common"
	ethtypes "github.com/ledgerwatch/erigon/core/types"
	"github.com/ledgerwatch/erigon/ethclient"
)

// var incorrectlyProcessTx int = 0

type TxTarget struct {
	targetRpcEndpoint string
	ethClient         *ethclient.Client
	faucetPrivateKey  string
	fundingAmount     uint64
	txSendingLimit    int64
}

func New(targetRpcEndpoint string, ethClient *ethclient.Client, faucetPrivateKey string, fundingAmount uint64, txSendingLimit int64) *TxTarget {
	return &TxTarget{
		targetRpcEndpoint: targetRpcEndpoint,
		ethClient:         ethClient,
		faucetPrivateKey:  faucetPrivateKey,
		fundingAmount:     fundingAmount,
		txSendingLimit:    txSendingLimit,
	}
}

func (tt *TxTarget) SendTxs(txs []*types.Tx) error {
	var err error

	txsCount := len(txs)
	startTime := time.Now()

	for i, tx := range txs {
		if _, err = tt.SendTx(tx.Hash, tx.Bytes); err != nil {
			return err
		}

		if (i+1)&(SENDING_QUEUE_SIZE-1) == 0 {
			fmt.Printf("Sent %d transactions\n", i+1)
			tt.waitForTxPoolToGoBelowLimit(SENDING_QUEUE_SIZE)
		}

		if tt.shouldApplyTxLimit() {
			timeInMicroAtTheEndOfThisTx := int64(i+1) * 1000000 / tt.txSendingLimit
			sleepTimeInMicro := timeInMicroAtTheEndOfThisTx - time.Since(startTime).Microseconds()
			if sleepTimeInMicro > 0 {
				time.Sleep(time.Duration(sleepTimeInMicro) * time.Microsecond)
			}
		}
	}

	fmt.Printf("Sent %d transactions for %.3f seconds\n", txsCount, float32(time.Since(startTime).Milliseconds())/1000)

	fmt.Println()
	fmt.Printf("Waiting for %d transactions\n", txsCount)

	// wait for tx pool to be empty
	tt.waitForEmptyTxPool()
	// wait for last 256 txs then assume that all of them has already been processed
	last256Index := txsCount - 256
	if last256Index < 0 {
		last256Index = 0
	}
	for i := txsCount - 1; i >= last256Index; i-- {
		tt.waitTx(txs[i])
	}

	timeIncludingStart := time.Since(startTime).Seconds()
	fmt.Printf("Executing transactions at %f tx/sec. rate\n", float64(txsCount)/timeIncludingStart)
	fmt.Println()

	fmt.Printf("Getting receipts of %d transactions\n", txsCount)
	totalGas := uint64(0)
	for i := txsCount - 1; i >= 0; i-- {
		receipt, err := tt.fetchReceipt(txs[i].Hash)
		if err != nil {
			return fmt.Errorf("tx %d was not processed but the tx pool was already empty: %v", i, err)
		}
		totalGas += receipt.GasUsed
	}

	// fmt.Println()
	// fmt.Printf("Incorrectly completed %d out of %d\n", incorrectlyProcessTx, txsCount)

	fmt.Printf("Gas %.2f gas/sec.\n", float64(totalGas)/timeIncludingStart)
	fmt.Println()

	return nil
}

func (tt *TxTarget) SendTx(txHash string, rlp []byte) (string, error) {
	hexEncodedTx := hex.EncodeToString(rlp)
	var resp []byte
	var err error

	if resp, err = utils.MakePostRequest(tt.targetRpcEndpoint, makeTxSendParams(hexEncodedTx)); err != nil {
		return "", err
	}

	var transactionRes types.SendRawTransactionResponse
	err = json.Unmarshal(resp, &transactionRes)
	if err != nil {
		return "", err
	}

	if transactionRes.Error != nil {
		return "", fmt.Errorf("hash (%s): %s", txHash, transactionRes.Error.Message)
	}

	return transactionRes.Result, nil
}

func (tt *TxTarget) waitTx(tx *types.Tx) *ethtypes.Receipt {
	for {
		receipt, err := tt.fetchReceipt(tx.Hash)
		if err != nil {
			time.Sleep(10 * time.Microsecond)
			continue
		}

		// if receipt.Status != tx.Status {
		// 	incorrectlyProcessTx++
		// }

		return receipt
	}
}

func (tt *TxTarget) isTxSuccessful(txHash string) (bool, error) {
	for {
		receipt, err := tt.fetchReceipt(txHash)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		return receipt.Status == 1, nil
	}
}

func (tt *TxTarget) fetchReceipt(txHash string) (*ethtypes.Receipt, error) {
	ctx := context.Background()
	return tt.ethClient.TransactionReceipt(ctx, common.HexToHash(txHash))
}
