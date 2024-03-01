package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gateway-fm/tx-repeater/src/persistor"
	txsource "github.com/gateway-fm/tx-repeater/src/tx_source"
	txtarget "github.com/gateway-fm/tx-repeater/src/tx_target"
	txtargettypes "github.com/gateway-fm/tx-repeater/src/tx_target/types"
)

func main() {
	var txs []*txtargettypes.Tx
	var err error

	var source string
	var target string
	var txCount int
	var faucetPrivateKey string
	var fundingAmount int64
	var txSendingLimit int64

	flag.StringVar(&source, "source", "https://zkevm-rpc.com", "RPC address to get transactions from")
	flag.StringVar(&target, "destination", "http://localhost:8467", "RPC addresses to send transactions to")
	flag.StringVar(&faucetPrivateKey, "faucet-key", "", "Private key of the faucet wallet")
	flag.IntVar(&txCount, "tx-count", 0, "Block number to start from")
	flag.Int64Var(&fundingAmount, "funding-amount", 200, "This indicates how many ETH each account will be pre-funded")
	flag.Int64Var(&txSendingLimit, "tx-sending-limit", -1, "Limit how many transactions per second are sent to the target")
	flag.Parse()

	ethClient, err := ethclient.Dial(target)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		return
	}
	defer ethClient.Close()

	currentWorkingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		return
	}

	persistor, err := persistor.New(currentWorkingDir + "/data")
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		return
	}

	txSource := txsource.New(source, persistor)
	if txs, err = txSource.FetchAllTransactions(txCount); err != nil {
		fmt.Printf("error: %+v\n", err)
		return
	}

	txTarget := txtarget.New(target, ethClient, faucetPrivateKey, fundingAmount, txSendingLimit)
	if err := txTarget.EnsureFunding(txs); err != nil {
		fmt.Printf("error: %+v\n", err)
		return
	}

	if err = txTarget.SendTxs(txs); err != nil {
		fmt.Printf("error: %+v\n", err)
		return
	}

	fmt.Printf("Done\n")
}
