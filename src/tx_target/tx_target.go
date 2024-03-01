package txtarget

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	types "github.com/gateway-fm/tx-repeater/src/tx_target/types"
	"github.com/gateway-fm/tx-repeater/src/utils"
)

var incorrectlyProcessTx int = 0

type TxTarget struct {
	endpoint         string
	ethClient        *ethclient.Client
	faucetPrivateKey string
	fundingAmount    int64
	txSendingLimit   int64
}

func New(endpoint string, ethClient *ethclient.Client, faucetPrivateKey string, fundingAmount, txSendingLimit int64) *TxTarget {
	return &TxTarget{
		endpoint:         endpoint,
		ethClient:        ethClient,
		faucetPrivateKey: faucetPrivateKey,
		fundingAmount:    fundingAmount,
		txSendingLimit:   txSendingLimit,
	}
}

func (tt *TxTarget) SendTxs(txs []*types.Tx) error {
	var err error

	startTime := time.Now()

	for i, tx := range txs {
		if _, err = tt.SendTx(tx.Hash, tx.Bytes); err != nil {
			return err
		}

		if tt.shouldApplyTxLimit() {
			timeInMicroAtTheEndOfThisTx := int64(i+1) * 1000000 / tt.txSendingLimit
			sleepTimeInMicro := timeInMicroAtTheEndOfThisTx - time.Since(startTime).Microseconds()
			if sleepTimeInMicro > 0 {
				time.Sleep(time.Duration(sleepTimeInMicro) * time.Microsecond)
			}
		}
	}

	if tt.shouldApplyTxLimit() {
		fmt.Printf("Sent %d transactions for %.3f seconds\n", len(txs), float32(time.Since(startTime).Milliseconds())/1000)
	}

	fmt.Printf("Waiting for %d transactions\n", len(txs))
	startTimeForPerfMeasurement := time.Now()
	for i := len(txs) - 1; i >= 0; i-- {
		if ok := tt.isTxProcessed(txs[i]); !ok {
			return fmt.Errorf("tx not processed: %s", txs[i].Hash)
		}
	}

	fmt.Printf("Incorrectly completed %d out of %d\n", incorrectlyProcessTx, len(txs))
	fmt.Printf("Executing transactions at %f tx/sec. rate (including SEND time)\n", float64(len(txs))/time.Since(startTime).Seconds())
	fmt.Printf("Executing transactions at %f tx/sec. rate (excluding SEND time)\n", float64(len(txs))/time.Since(startTimeForPerfMeasurement).Seconds())

	return nil
}

func (tt *TxTarget) SendTx(txHash string, rlp []byte) (string, error) {
	hexEncodedTx := hex.EncodeToString(rlp)
	var resp []byte
	var err error

	if resp, err = utils.MakePostRequest(tt.endpoint, makeTxSendParams(hexEncodedTx)); err != nil {
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

func (tt *TxTarget) isTxProcessed(tx *types.Tx) bool {
	ctx := context.Background()

	for {
		receipt, err := tt.ethClient.TransactionReceipt(ctx, common.HexToHash(tx.Hash))
		if err != nil {
			time.Sleep(10 * time.Microsecond)
			continue
		}

		if receipt.Status != tx.Status {
			incorrectlyProcessTx++
		}

		return true
	}
}

func (tt *TxTarget) isTxSuccessful(txHash string) (bool, error) {
	ctx := context.Background()

	for {
		receipt, err := tt.ethClient.TransactionReceipt(ctx, common.HexToHash(txHash))
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		return receipt.Status == 1, nil
	}
}

func (tt *TxTarget) getBalanceByHexAddress(address string) (*big.Int, error) {
	return tt.ethClient.BalanceAt(context.Background(), common.HexToAddress(address), nil)
}

func (tt *TxTarget) shouldApplyTxLimit() bool {
	return tt.txSendingLimit > 0
}
