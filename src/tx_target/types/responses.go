package types

import (
	"fmt"
	"strconv"
)

type SendRawTransactionResponse struct {
	Result string         `json:"result"`
	Error  *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type TxPoolStatusResponse struct {
	Result *TxPoolStatusResultResponse `json:"result"`
}

func (t *TxPoolStatusResponse) IsTxPoolEmpty() bool {
	return t.Result.getBaseFeeCount() == 0 && t.Result.getPendingCount() == 0 && t.Result.getQueuedCount() == 0
}

func (t *TxPoolStatusResponse) IsTxPoolLessThan(limit uint64) bool {
	return t.Result.getBaseFeeCount()+t.Result.getPendingCount()+t.Result.getQueuedCount() < limit
}

type TxPoolStatusResultResponse struct {
	BaseFee string `json:"baseFee"`
	Pending string `json:"pending"`
	Queued  string `json:"queued"`
}

func (t *TxPoolStatusResultResponse) getBaseFeeCount() uint64 {
	result, err := strconv.ParseUint(t.BaseFee[2:], 16, 64)
	if err != nil {
		panic(fmt.Errorf("cannot parse %s: %v", t.BaseFee[2:], err))
	}
	return result
}

func (t *TxPoolStatusResultResponse) getPendingCount() uint64 {
	result, err := strconv.ParseUint(t.Pending[2:], 16, 64)
	if err != nil {
		panic(fmt.Errorf("cannot parse %s: %v", t.Pending[2:], err))
	}
	return result
}

func (t *TxPoolStatusResultResponse) getQueuedCount() uint64 {
	result, err := strconv.ParseUint(t.Queued[2:], 16, 64)
	if err != nil {
		panic(fmt.Errorf("cannot parse %s: %v", t.Queued[2:], err))
	}
	return result
}
