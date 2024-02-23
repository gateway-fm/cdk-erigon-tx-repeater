package txsource

import (
	"bytes"
	"fmt"
	"io"
)

func makeBlockReqParams(blockNumber int) io.Reader {
	paramsAsString := fmt.Sprintf(`
	{
		"jsonrpc":"2.0",
		"method":"eth_getBlockByNumber",
		"params":[
			"0x%x", 
			false
		],
		"id":"1"
	}
	`, blockNumber)

	return bytes.NewReader([]byte(paramsAsString))
}

func makeTransactionReqParams(txHash string) io.Reader {
	paramsAsString := fmt.Sprintf(`
	{
		"jsonrpc":"2.0",
		"method":"eth_getTransactionByHash",
		"params":[
			"%s"
		],
		"id":"1"
	}
	`, txHash)

	return bytes.NewReader([]byte(paramsAsString))
}
