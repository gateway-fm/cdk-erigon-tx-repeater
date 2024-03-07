package txtarget

import (
	"encoding/json"
	"time"

	types "github.com/gateway-fm/tx-repeater/src/tx_target/types"
	"github.com/gateway-fm/tx-repeater/src/utils"
)

const SENDING_QUEUE_SIZE = 1 << 12 // MUST BE POWER OF 2

func (tt *TxTarget) fetchTxPoolStatus() *types.TxPoolStatusResponse {
	var resp []byte
	var err error

	for {
		if resp, err = utils.MakePostRequest(tt.targetRpcEndpoint, makeTxPoolStatus()); err != nil {
			time.Sleep(64 * time.Millisecond)
			continue
		}

		var txPoolStatusResp types.TxPoolStatusResponse
		err = json.Unmarshal(resp, &txPoolStatusResp)
		if err != nil {
			time.Sleep(64 * time.Millisecond)
			continue
		}

		return &txPoolStatusResp
	}
}

func (tt *TxTarget) waitForEmptyTxPool() bool {
	for {
		txPoolStatusResp := tt.fetchTxPoolStatus()

		if !txPoolStatusResp.IsTxPoolEmpty() {
			time.Sleep(64 * time.Millisecond)
			continue
		}

		return true
	}
}

func (tt *TxTarget) waitForTxPoolToGoBelowLimit(limit uint64) bool {
	for {
		txPoolStatusResp := tt.fetchTxPoolStatus()

		if !txPoolStatusResp.IsTxPoolLessThan(limit) {
			time.Sleep(64 * time.Millisecond)
			continue
		}

		return true
	}
}
