package main

import (
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

	ethClient, err := ethclient.Dial("http://localhost:8467")
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

	txSource := txsource.New("https://zkevm-rpc.com", persistor)
	if txs, err = txSource.FetchAllTransactions(128); err != nil {
		fmt.Printf("error: %+v\n", err)
		return
	}

	txTarget := txtarget.New("http://localhost:8467", ethClient)
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
