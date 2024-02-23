package utils

import (
	"io"
	"net/http"
)

const CHAIN_ID = 1101
const BRIDGE_ADDRESS = "0x2a3dd3eb832af982ec71669e178424b10dca2ede"

func MakePostRequest(endpoint string, params io.Reader) ([]byte, error) {
	httpReq, err := http.NewRequest("POST", endpoint, params)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}
