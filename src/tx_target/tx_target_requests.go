package txtarget

import (
	"bytes"
	"fmt"
	"io"
)

func makeTxSendParams(hexEncodedTx string) io.Reader {
	paramsAsString := fmt.Sprintf(`
	{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0x%s"],"id":1}
	`, hexEncodedTx)

	return bytes.NewReader([]byte(paramsAsString))
}

func makeTxPoolStatus() io.Reader {
	return bytes.NewReader([]byte(`{"jsonrpc":"2.0","method":"txpool_status","params":[],"id":1}`))
}
