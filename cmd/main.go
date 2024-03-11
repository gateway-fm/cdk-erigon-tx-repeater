package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gateway-fm/tx-repeater/src/persistor"
	txsource "github.com/gateway-fm/tx-repeater/src/tx_source"
	txtarget "github.com/gateway-fm/tx-repeater/src/tx_target"
	txtargettypes "github.com/gateway-fm/tx-repeater/src/tx_target/types"
	"github.com/gateway-fm/tx-repeater/src/utils"
	"github.com/ledgerwatch/erigon/ethclient"
)

func main() {
	var txs []*txtargettypes.Tx
	var err error

	var flagSourceDatastreamEndpoint string
	var flagSourceRpcEndpoint string
	var flagTargetRpcEndpoint string
	var flagFaucetPrivateKey string
	var flagBlocksCount uint64
	var flagFundingAmount uint64
	var flagTxSendingLimit int64

	flag.StringVar(&flagSourceDatastreamEndpoint, "source-datastream-endpoint", "stream.zkevm-rpc.com:6900", "Source datastream URL")
	flag.StringVar(&flagSourceRpcEndpoint, "source-rpc-endpoint", "stream.zkevm-rpc.com:6900", "Source RPC URL")
	flag.StringVar(&flagTargetRpcEndpoint, "target-rpc-endpoint", "http://localhost:8467", "RPC URL to send transactions to")
	flag.StringVar(&flagFaucetPrivateKey, "faucet-key", "", "Private key of the faucet wallet")
	flag.Uint64Var(&flagBlocksCount, "blocks", 0, "Number of blocks to fetch from the source")
	flag.Uint64Var(&flagFundingAmount, "funding-amount", 200, "This indicates how many ETH each account will be pre-funded")
	flag.Int64Var(&flagTxSendingLimit, "tx-sending-limit", -1, "Limit how many transactions per second are sent to the target")
	flag.Int64Var(&utils.CHAIN_ID, "chain-id", 1101, "The chain-id of the sequencer")
	flag.Parse()

	ethClient, err := ethclient.Dial(flagTargetRpcEndpoint)
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

	txSource := txsource.New(flagSourceDatastreamEndpoint, flagSourceRpcEndpoint, persistor)
	if txs, err = txSource.FetchAllTransactions(flagBlocksCount); err != nil {
		fmt.Printf("error: %+v\n", err)
		return
	}

	txTarget := txtarget.New(flagTargetRpcEndpoint, ethClient, flagFaucetPrivateKey, flagFundingAmount, flagTxSendingLimit)
	if err := txTarget.EnsureFunding(txs); err != nil {
		fmt.Printf("error: %+v\n", err)
		return
	}

	if err = txTarget.SendTxs(txs); err != nil {
		fmt.Printf("error: %+v\n", err)
		return
	}

	fmt.Printf("\nDone\n")
}
