package txsource

import (
	"encoding/json"

	persistor "github.com/gateway-fm/tx-repeater/src/persistor"
	types "github.com/gateway-fm/tx-repeater/src/tx_source/types"
	txtargettypes "github.com/gateway-fm/tx-repeater/src/tx_target/types"
	utils "github.com/gateway-fm/tx-repeater/src/utils"
)

type TxSource struct {
	endpoint  string
	persistor *persistor.Persistor
}

func New(endpoint string, persistor *persistor.Persistor) *TxSource {
	return &TxSource{
		endpoint:  endpoint,
		persistor: persistor,
	}
}

func (ts *TxSource) FetchAllTransactions(minNumberOfTx int) ([]*txtargettypes.Tx, error) {
	var block *txtargettypes.Block
	var err error

	txs := []*txtargettypes.Tx{}
	latestBlock := 1

	if ts.persistor != nil {
		latestBlock = ts.persistor.FetchLatestBlockNumber()
		for i := 1; i <= latestBlock; i++ {
			if block, err = ts.persistor.FetchBlock(i); err != nil {
				return nil, err
			}

			if block != nil {
				txs = append(txs, block.Transactions...)
			}
		}
		latestBlock++
	}

	for {
		if latestBlock >= 3000000 {
			break
		}
		if len(txs) >= minNumberOfTx {
			break
		}

		if block, err = ts.fetchBlock(latestBlock); err != nil {
			return nil, err
		}

		if ts.persistor != nil {
			if block.HasTransactions() {
				if err := ts.persistor.CreditBlock(block); err != nil {
					return nil, err
				}
			}
			if err := ts.persistor.CreditLatestBlockNumber(block.Number); err != nil {
				return nil, err
			}
		}

		txs = append(txs, block.Transactions...)
		latestBlock++
	}

	return txs, nil
}

func (ts *TxSource) fetchBlock(blockNumber int) (*txtargettypes.Block, error) {
	var txTarget *txtargettypes.Tx
	var resp []byte
	var err error

	if resp, err = utils.MakePostRequest(ts.endpoint, makeBlockReqParams(blockNumber)); err != nil {
		return nil, err
	}

	var blockSource types.Block
	if err = json.Unmarshal(resp, &blockSource); err != nil {
		return nil, err
	}

	block := txtargettypes.NewBlock(blockNumber, len(blockSource.Result.Transactions))

	for _, txHash := range blockSource.Result.Transactions {
		if txTarget, err = ts.fetchTransaction(txHash); err != nil {
			return nil, err
		}
		if txTarget != nil {
			block.AppendTx(txTarget)
		}
	}

	return block, nil
}

func (ts *TxSource) fetchTransaction(txHash string) (*txtargettypes.Tx, error) {
	var resp []byte
	var err error

	if resp, err = utils.MakePostRequest(ts.endpoint, makeTransactionReqParams(txHash)); err != nil {
		return nil, err
	}

	var txSource types.Tx
	if err = json.Unmarshal(resp, &txSource); err != nil {
		return nil, err
	}

	// if txSource.Result.V == "0x1b" || txSource.Result.V == "0x1c" || txSource.Result.V == "0x0" || txSource.Result.V == "0x1" {
	// 	fmt.Printf("Drop tx %s\n", txHash)
	// 	return nil, nil
	// }

	if txSource.IsBridgeTx() {
		return nil, nil
	}

	return txtargettypes.FromSourceTx(&txSource)

}
