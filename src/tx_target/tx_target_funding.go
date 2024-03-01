package txtarget

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	types "github.com/gateway-fm/tx-repeater/src/tx_target/types"
	"github.com/gateway-fm/tx-repeater/src/utils"
)

func (tt *TxTarget) EnsureFunding(txs []*types.Tx) error {
	var balance *big.Int
	var fundTx []byte
	var err error
	var txHash string
	fromsMap := map[string]string{}

	for _, tx := range txs {
		fromsMap[tx.From] = tx.From
	}

	fmt.Printf("Checking %d transactions for initial funding\n", len(fromsMap))
	fundTxs := make([][]byte, 0, len(fromsMap))
	nonceAddition := uint64(0)
	for from := range fromsMap {
		if balance, err = tt.getBalanceByHexAddress(from); err != nil {
			return err
		}

		if balance.Cmp(big.NewInt(2000000000000000000)) >= 0 {
			continue
		}

		if fundTx, err = tt.makeFundingTx(from, nonceAddition); err != nil {
			return err
		}
		fundTxs = append(fundTxs, fundTx)
		nonceAddition++
	}

	fmt.Printf("Broadcasting %d transactions for preparing initial funding\n", len(fundTxs))
	fundTxsHashes := make([]string, 0, len(fundTxs))
	for _, fundTx := range fundTxs {
		if txHash, err = tt.SendTx("", fundTx); err != nil {
			return err
		}
		fundTxsHashes = append(fundTxsHashes, txHash)
	}

	for _, txHash := range fundTxsHashes {
		isSuccessful, err := tt.isTxSuccessful(txHash)
		if err != nil {
			return err
		}
		if !isSuccessful {
			return fmt.Errorf("tx %s was not successful", txHash)
		}
	}

	fmt.Printf("Finish initial funding\n")

	return nil
}

func (tt *TxTarget) makeFundingTx(recipientAddress string, nonceAddition uint64) ([]byte, error) {
	privateKey, err := crypto.HexToECDSA(tt.faucetPrivateKey)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := tt.ethClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	toAddress := common.HexToAddress(recipientAddress)
	txData := &ethtypes.LegacyTx{
		Nonce:    nonce + nonceAddition,
		To:       &toAddress,
		Value:    big.NewInt(0).Mul(big.NewInt(1000000000000000000), big.NewInt(tt.fundingAmount)),
		Gas:      uint64(1000000),
		GasPrice: big.NewInt(100000000000),
		Data:     nil,
		V:        big.NewInt(2237),
	}

	signedTx, err := ethtypes.SignNewTx(privateKey, ethtypes.LatestSignerForChainID(big.NewInt(utils.CHAIN_ID)), txData)
	if err != nil {
		return nil, err
	}

	var rawTx bytes.Buffer
	err = signedTx.EncodeRLP(&rawTx)
	if err != nil {
		return nil, err
	}

	return rawTx.Bytes(), nil

}
