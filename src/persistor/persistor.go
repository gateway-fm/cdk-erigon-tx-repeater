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

const BLOCKS_PER_FILE = 1 << 10 // MUST BE POWER OF 2
const BLOCKS_PER_FILE_ = BLOCKS_PER_FILE - 1

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

func (p *Persistor) FetchLatestBlockNumber() uint64 {
	blockNumberContentAsBytes, err := os.ReadFile(p.makeFilenameForLatestBlock())
	if err != nil {
		return 0
	}

	blockNumberContentAsString := strings.Trim(string(blockNumberContentAsBytes), "\n\r")

	blockNum, err := strconv.ParseUint(blockNumberContentAsString, 10, 64)
	if err != nil {
		return 0
	}

	return blockNum
}

func (p *Persistor) CreditLatestBlockNumber(blockNum uint64) error {
	blockNumberContentAsString := strconv.FormatUint(blockNum, 10)
	blockNumberContentAsBytes := []byte(blockNumberContentAsString)

	return os.WriteFile(p.makeFilenameForLatestBlock(), blockNumberContentAsBytes, 0755)
}

func (p *Persistor) FetchBlocks(startBlockNum uint64, endBlockNum uint64) ([]*txtargettypes.Block, error) {
	var lastReadFilename string
	var blocksPerFile []*txtargettypes.Block
	var err error

	blocks := make([]*txtargettypes.Block, 0, endBlockNum-startBlockNum+1)

	for blockNum := startBlockNum; blockNum <= endBlockNum; blockNum++ {
		blockFilename := p.makeFilenameByBlockNum(blockNum)
		if lastReadFilename != blockFilename {
			lastReadFilename = blockFilename
			if blocksPerFile, err = p.fetchBlocksFromFile(blockFilename); err != nil {
				return nil, err
			}
		}

		block := blocksPerFile[blockNum&BLOCKS_PER_FILE_]
		if block != nil {
			blocks = append(blocks, block)
		}
	}

	return blocks, nil
}

func (p *Persistor) fetchBlocksFromFile(filename string) ([]*txtargettypes.Block, error) {
	blockContentAsBytes, err := os.ReadFile(filename)
	if err != nil {
		return make([]*txtargettypes.Block, BLOCKS_PER_FILE), nil
	}

	var blocks []*txtargettypes.Block
	err = json.Unmarshal(blockContentAsBytes, &blocks)
	if err != nil {
		return nil, err
	}

	return blocks, nil
}

func (p *Persistor) CreditBlocks(blocks []*txtargettypes.Block) error {
	var lastReadFilename string
	var blocksPerFile []*txtargettypes.Block
	var err error

	for _, block := range blocks {
		blockFilename := p.makeFilenameByBlockNum(block.Number)
		if lastReadFilename != blockFilename {
			//save the content so far
			if lastReadFilename != "" {
				p.creditBlocksFromFile(lastReadFilename, blocksPerFile)
			}
			//load new content
			lastReadFilename = blockFilename
			if blocksPerFile, err = p.fetchBlocksFromFile(blockFilename); err != nil {
				return err
			}
		}

		blocksPerFile[block.Number&BLOCKS_PER_FILE_] = block
	}

	if lastReadFilename != "" {
		p.creditBlocksFromFile(lastReadFilename, blocksPerFile)
	}

	return nil
}

func (p *Persistor) creditBlocksFromFile(filename string, blocks []*txtargettypes.Block) error {
	blocksContentAsBytes, err := json.Marshal(blocks)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, blocksContentAsBytes, 0755)
}

func (p *Persistor) makeFilenameForLatestBlock() string {
	return p.path + "/blockNumber.dt"
}

func (p *Persistor) makeFilenameByBlockNum(blockNum uint64) string {
	lowestBlockNum := (blockNum >> 10) << 10
	highestBlockNum := lowestBlockNum | BLOCKS_PER_FILE_

	return fmt.Sprintf("%s/block%010d-%010d.dt", p.path, lowestBlockNum, highestBlockNum)
}
