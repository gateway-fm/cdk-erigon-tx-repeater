package utils

import (
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

var CHAIN_ID = int64(1101)

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

func StringToInt(input string) (int64, error) {
	return strconv.ParseInt(TrimHex(input), 16, 64)
}

func StringToUint(input string) (uint64, error) {
	return strconv.ParseUint(TrimHex(input), 16, 64)
}

func StringToBig(input string) *big.Int {
	result := new(big.Int)
	result.SetString(TrimHex(input), 16)
	return result
}

func TrimHex(input string) string {
	if strings.HasPrefix(input, "0x") {
		return input[2:]
	}
	return input
}
