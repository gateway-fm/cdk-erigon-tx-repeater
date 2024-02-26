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

type TxTarget struct {
	endpoint         string
	ethClient        *ethclient.Client
	faucetPrivateKey string
}

func New(endpoint string, ethClient *ethclient.Client, faucetPrivateKey string) *TxTarget {
	return &TxTarget{
		endpoint:         endpoint,
		ethClient:        ethClient,
		faucetPrivateKey: faucetPrivateKey,
	}
}

func (tt *TxTarget) SendTxs(txs []*types.Tx) error {
	var txHash string
	var err error

	txHashes := make([]string, 0, len(txs))
	fromsMap := map[string]bool{}

	for i, tx := range txs {
		if _, ok := fromsMap[tx.From]; ok {
			fmt.Printf("Waiting for batch with %d entries\n", len(txHashes))
			if err := tt.areTxsProcessed(txHashes); err != nil {
				return nil
			}

			fromsMap = make(map[string]bool)
			txHashes = make([]string, 0, len(txs))
		}

		if txHash, err = tt.SendTx(tx.Bytes); err != nil {
			return err
		}
		fromsMap[tx.From] = true
		txHashes = append(txHashes, txHash)
		fmt.Printf("%d: %s\n", i, txHash)
	}

	if err := tt.areTxsProcessed(txHashes); err != nil {
		return nil
	}

	return nil
}

func (tt *TxTarget) SendTx(rlp []byte) (string, error) {
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
		return "", fmt.Errorf("%s", transactionRes.Error.Message)
	}

	return transactionRes.Result, nil
}

func (tt *TxTarget) areTxsProcessed(txHashes []string) error {
	for _, txHash := range txHashes {
		if ok := tt.isTxProcessed(txHash); !ok {
			return fmt.Errorf("tx not processed: %s", txHash)
		}
	}

	return nil
}

func (tt *TxTarget) isTxProcessed(txHash string) bool {
	ctx := context.Background()

	for {
		_, err := tt.ethClient.TransactionReceipt(ctx, common.HexToHash(txHash))
		if err != nil {
			time.Sleep(10 * time.Microsecond)
			continue
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

func (tt *TxTarget) printBalance(address string) {
	balance, err := tt.getBalanceByHexAddress(address)
	if err != nil {
		fmt.Println("Error getting balance")
	}
	fmt.Printf("Balance of %s [%v] = %s\n", address, common.HexToAddress(address), balance.String())
}
