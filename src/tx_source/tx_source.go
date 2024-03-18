package txsource

import (
	"fmt"
	"math/big"

	persistor "github.com/gateway-fm/tx-repeater/src/persistor"
	txtargettypes "github.com/gateway-fm/tx-repeater/src/tx_target/types"
	"github.com/gateway-fm/tx-repeater/src/utils"
	"github.com/holiman/uint256"
	"github.com/ledgerwatch/erigon-lib/common"
	ethtypes "github.com/ledgerwatch/erigon/core/types"
	"github.com/ledgerwatch/erigon/zk/datastream"
	txtype "github.com/ledgerwatch/erigon/zk/tx"
)

type TxSource struct {
	datastreamEndpoint string
	rpcEndpoint        string
	persistor          *persistor.Persistor
}

func New(datastreamEndpoint, rpcEndpoint string, persistor *persistor.Persistor) *TxSource {
	return &TxSource{
		datastreamEndpoint: datastreamEndpoint,
		rpcEndpoint:        rpcEndpoint,
		persistor:          persistor,
	}
}

func (ts *TxSource) FetchAllTransactions(targetBlocksCount uint64) ([]*txtargettypes.Tx, error) {
	var blocks []*txtargettypes.Block
	var err error

	txs := []*txtargettypes.Tx{}
	currentBlock := uint64(1)

	if ts.persistor != nil {
		persistedBlocksCount := ts.persistor.FetchLatestBlockNumber()
		endBlockNumber := targetBlocksCount
		if endBlockNumber > persistedBlocksCount {
			endBlockNumber = persistedBlocksCount
		}

		fmt.Printf("Reading %d blocks from disk\n", endBlockNumber)
		if blocks, err = ts.persistor.FetchBlocks(1, endBlockNumber); err != nil {
			return nil, err
		}

		for _, block := range blocks {
			txs = append(txs, block.Transactions...)
		}
		currentBlock = persistedBlocksCount + 1
	}

	for {
		if currentBlock > targetBlocksCount {
			break
		}

		requestedBlocksCount := uint64(4096)
		maxBlocksCount := targetBlocksCount - (currentBlock - 1)
		if requestedBlocksCount > maxBlocksCount {
			requestedBlocksCount = maxBlocksCount
		}

		fmt.Printf("Reading [%d-%d] blocks from datastream\n", currentBlock, currentBlock+requestedBlocksCount-1)
		if blocks, err = ts.fetchBlocks(currentBlock, int(requestedBlocksCount)); err != nil {
			return nil, err
		}

		fetchedBlocksCount := uint64(len(blocks))

		if ts.persistor != nil {
			filteredBlocks := make([]*txtargettypes.Block, 0, fetchedBlocksCount)
			for _, block := range blocks {
				if block.HasTransactions() {
					filteredBlocks = append(filteredBlocks, block)
				}
			}

			if err = ts.persistor.CreditBlocks(filteredBlocks); err != nil {
				return nil, err
			}
			if fetchedBlocksCount > 0 {
				if err := ts.persistor.CreditLatestBlockNumber(blocks[fetchedBlocksCount-1].Number); err != nil {
					return nil, err
				}
			}
		}

		for _, block := range blocks {
			txs = append(txs, block.Transactions...)
		}

		currentBlock += fetchedBlocksCount
	}

	return txs, nil
}

func (ts *TxSource) fetchBlocks(startBlockNumber uint64, blocksCount int) ([]*txtargettypes.Block, error) {
	l2blocks, _, _, _, err := datastream.DownloadL2Blocks(ts.datastreamEndpoint, startBlockNumber, blocksCount)
	if err != nil {
		return nil, err
	}

	var txTo *string
	var txToAddr *common.Address
	var blocks []*txtargettypes.Block = make([]*txtargettypes.Block, 0, len(*l2blocks))

	signer := *(ethtypes.LatestSignerForChainID(big.NewInt(utils.CHAIN_ID)))

	for _, l2block := range *l2blocks {
		block := txtargettypes.NewBlock(l2block.L2BlockNumber, len(l2block.L2Txs))
		blocks = append(blocks, block)

		for _, tx := range l2block.L2Txs {
			ltx, _, err := txtype.DecodeTx(tx.Encoded, tx.EffectiveGasPricePercentage, (*l2blocks)[0].ForkId)
			if err != nil {
				return nil, err
			}

			from, err := ltx.Sender(signer)
			if err != nil {
				return nil, err
			}

			txToAddr = ltx.GetTo()
			if txToAddr != nil {
				txToValue := txToAddr.Hex()
				txTo = &txToValue
			} else {
				txTo = nil
			}

			txFunds := uint256.NewInt(500000000000000000)

			txCost := ltx.Cost()
			if txCost != nil {
				txFunds = txFunds.Add(txFunds, txCost)
			}

			txValue := ltx.GetValue()
			if txValue != nil {
				txFunds = txFunds.Add(txFunds, txValue)
			}

			txTarget := txtargettypes.NewTx(tx.Encoded, from.Hex(), txTo, ltx.Hash().Hex(), txFunds.Hex())
			block.AppendTx(txTarget)
		}
	}

	return blocks, nil
}
