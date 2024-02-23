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
	endpoint  string
	ethClient *ethclient.Client
}

func New(endpoint string, ethClient *ethclient.Client) *TxTarget {
	return &TxTarget{
		endpoint:  endpoint,
		ethClient: ethClient,
	}
}

func (tt *TxTarget) SendTxs(txs []*types.Tx) error {
	var txHash string
	var err error

	for i, tx := range txs {
		if txHash, err = tt.SendTx(tx.Bytes); err != nil {
			return err
		}
		fmt.Printf("%d: %s\n", i, txHash)
		if ok := tt.isTxProcessed(txHash); !ok {
			return fmt.Errorf("tx not processed: %s", txHash)
		}
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

func (tt *TxTarget) isTxProcessed(txHash string) bool {
	ctx := context.Background()

	for {
		_, err := tt.ethClient.TransactionReceipt(ctx, common.HexToHash(txHash))
		if err != nil {
			time.Sleep(1 * time.Second)
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
