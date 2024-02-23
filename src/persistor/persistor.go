package persistor

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	txtargettypes "github.com/gateway-fm/tx-repeater/src/tx_target/types"
)

type Persistor struct {
	path string
}

func New(path string) (*Persistor, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return nil, err
		}
	}

	return &Persistor{
		path: path,
	}, nil
}

func (p *Persistor) FetchLatestBlockNumber() int {
	blockNumberContentAsBytes, err := os.ReadFile(p.makeFilenameForLatestBlock())
	if err != nil {
		return 0
	}

	blockNumberContentAsString := strings.Trim(string(blockNumberContentAsBytes), "\n\r")

	blockNum, err := strconv.Atoi(blockNumberContentAsString)
	if err != nil {
		return 0
	}

	return blockNum
}

func (p *Persistor) CreditLatestBlockNumber(blockNum int) error {
	blockNumberContentAsString := strconv.Itoa(blockNum)
	blockNumberContentAsBytes := []byte(blockNumberContentAsString)

	return os.WriteFile(p.makeFilenameForLatestBlock(), blockNumberContentAsBytes, 0755)
}

func (p *Persistor) FetchBlock(blockNum int) (*txtargettypes.Block, error) {
	blockContentAsBytes, err := os.ReadFile(p.makeFilenameByBlockNum(blockNum))
	if err != nil {
		return nil, nil
	}

	var block txtargettypes.Block
	err = json.Unmarshal(blockContentAsBytes, &block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func (p *Persistor) CreditBlock(block *txtargettypes.Block) error {
	blockContentAsBytes, err := json.Marshal(*block)
	if err != nil {
		return err
	}

	return os.WriteFile(p.makeFilenameByBlockNum(block.Number), blockContentAsBytes, 0755)
}

func (p *Persistor) makeFilenameForLatestBlock() string {
	return p.path + "/blockNumber.dt"
}

func (p *Persistor) makeFilenameByBlockNum(blockNum int) string {
	return fmt.Sprintf("%s/block%010d.dt", p.path, blockNum)
}
