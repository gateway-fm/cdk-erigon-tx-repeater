package types

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	txsourcetypes "github.com/gateway-fm/tx-repeater/src/tx_source/types"
	"github.com/umbracle/fastrlp"
)

type Tx struct {
	Bytes []byte `json:"bytes"`
	From  string `json:"from"`
	Hash  string `json:"hash"`
}

type TxType int64

const (
	LegacyTx     TxType = 0x0
	AccessListTx TxType = 0x01
	StateTx      TxType = 0x7f
	DynamicFeeTx TxType = 0x02
)

func NewTx(bytes []byte, from, hash string) *Tx {
	return &Tx{
		Bytes: bytes,
		From:  from,
		Hash:  hash,
	}
}

func FromSourceTx(resp *txsourcetypes.Tx) (*Tx, error) {
	// arena := fastrlp.DefaultArenaPool.Get()
	// defer fastrlp.DefaultArenaPool.Put(arena)

	// vv := arena.NewArray()

	t := resp.Result

	// typInt, err := stringToInt(t.Type)
	// if err != nil {
	// 	return nil, err
	// }
	// typ := TxType(typInt)

	// chainId := stringToBig(t.ChainID)
	gasPrice := stringToBig(t.GasPrice)
	value := stringToBig(t.Value)
	v := stringToBig(t.V)
	r := stringToBig(t.R)
	s := stringToBig(t.S)

	gas, err := stringToUint(t.Gas)
	if err != nil {
		return nil, err
	}

	// if gas == 0 {
	// 	gas = uint64(1000000)
	// }

	nonce, err := stringToUint(t.Nonce)
	if err != nil {
		return nil, err
	}

	// from, err := hex.DecodeString(trimHex(t.From))
	// if err != nil {
	// 	return nil, err
	// }

	var to []byte
	if t.To != nil {
		to, err = hex.DecodeString(trimHex(*t.To))
		if err != nil {
			return nil, err
		}
	}

	input, err := hex.DecodeString(trimHex(t.Input))
	if err != nil {
		return nil, err
	}

	toAddr := common.HexToAddress(hex.EncodeToString(to))
	toAddrPointer := &toAddr
	if t.To == nil {
		toAddrPointer = nil
	}
	txData := &ethtypes.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gas,
		To:       toAddrPointer,
		Value:    value,
		Data:     input,
		V:        v,
		R:        r,
		S:        s,
	}
	tx := ethtypes.NewTx(txData)
	b, _ := rlp.EncodeToBytes(tx)
	return NewTx(b, t.From, t.Hash), nil

	// if typ == AccessListTx {
	// 	vv.Set(arena.NewBigInt(chainId))
	// 	vv.Set(arena.NewUint(nonce))
	// 	vv.Set(arena.NewBigInt(gasPrice))
	// 	vv.Set(arena.NewUint(gas))
	// 	if t.To != nil {
	// 		vv.Set(arena.NewCopyBytes(to))
	// 	} else {
	// 		vv.Set(arena.NewNull())
	// 	}
	// 	vv.Set(arena.NewBigInt(value))
	// 	vv.Set(arena.NewCopyBytes(input))

	// 	err = RlpEncodeAccessList(arena, vv, t.AccessList)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	vv.Set(arena.NewBigInt(v))
	// 	vv.Set(arena.NewBigInt(r))
	// 	vv.Set(arena.NewBigInt(s))

	// 	dst := vv.MarshalTo(nil)
	// 	return NewTx(dst, t.From, t.Hash), nil
	// }

	// // Specify zero chain ID as per spec.
	// // This is needed to have the same format as other EVM chains do.
	// // There is no chain ID in the TX object, so it is always 0 here just to be compatible.
	// // Check Transaction1559Payload there https://eips.ethereum.org/EIPS/eip-1559#specification
	// if typ == DynamicFeeTx {
	// 	vv.Set(arena.NewBigInt(chainId))
	// }

	// vv.Set(arena.NewUint(nonce))

	// if typ == DynamicFeeTx {
	// 	tip := stringToBig(t.GasTipCap)
	// 	fee := stringToBig(t.GasFeeCap)

	// 	// Add EIP-1559 related fields.
	// 	// For non-dynamic-fee-tx gas price is used.
	// 	vv.Set(arena.NewBigInt(tip))
	// 	vv.Set(arena.NewBigInt(fee))
	// } else {
	// 	vv.Set(arena.NewBigInt(gasPrice))
	// }

	// vv.Set(arena.NewUint(gas))

	// // Address may be empty
	// if t.To != nil {
	// 	vv.Set(arena.NewCopyBytes(to))
	// } else {
	// 	vv.Set(arena.NewNull())
	// }

	// vv.Set(arena.NewBigInt(value))
	// vv.Set(arena.NewCopyBytes(input))

	// // Specify access list as per spec.
	// // This is needed to have the same format as other EVM chains do.
	// // There is no access list feature here, so it is always empty just to be compatible.
	// // Check Transaction1559Payload there https://eips.ethereum.org/EIPS/eip-1559#specification
	// if typ == DynamicFeeTx {
	// 	vv.Set(arena.NewArray())
	// }

	// // signature values
	// vv.Set(arena.NewBigInt(v))
	// vv.Set(arena.NewBigInt(r))
	// vv.Set(arena.NewBigInt(s))

	// if typ == StateTx {
	// 	vv.Set(arena.NewCopyBytes(from))
	// }

	// dst := vv.MarshalTo(nil)
	// return NewTx(dst, t.From, t.Hash), nil
}

func RlpEncodeAccessList(arena *fastrlp.Arena, vv *fastrlp.Value, list []txsourcetypes.TxAccessTuple) error {
	if len(list) == 0 {
		vv.Set(arena.NewNullArray())
	} else {
		ar1 := arena.NewArray()
		for _, at := range list {

			addr, err := hex.DecodeString(trimHex(at.Address))
			if err != nil {
				return err
			}

			ar2 := arena.NewArray()
			ar2.Set(arena.NewCopyBytes(addr))

			ar3 := arena.NewArray()
			for _, sKey := range at.StorageKeys {
				sk, err := hex.DecodeString(trimHex(sKey))
				if err != nil {
					return err
				}
				ar3.Set(arena.NewCopyBytes(sk))
			}
			ar2.Set(ar3)
		}
		vv.Set(ar1)
	}

	return nil
}

func stringToInt(input string) (int64, error) {
	return strconv.ParseInt(trimHex(input), 16, 64)
}

func stringToUint(input string) (uint64, error) {
	return strconv.ParseUint(trimHex(input), 16, 64)
}

func stringToBig(input string) *big.Int {
	result := new(big.Int)
	result.SetString(trimHex(input), 16)
	return result
}

func trimHex(input string) string {
	if strings.HasPrefix(input, "0x") {
		return input[2:]
	}
	return input
}
