package main

// import (
// 	"encoding/hex"
// 	"fmt"
// 	"math/big"
// 	"strings"

// 	"github.com/ethereum/go-ethereum/common"
// 	ethtypes "github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/crypto"
// 	"github.com/ethereum/go-ethereum/rlp"
// )

// func main() {
// 	nonce := uint64(0)
// 	input, _ := hex.DecodeString("2cffd02eb7cd745b9fc33c6e233768f51f262865c8cdff188d4e63c16709e389c11d5cd8ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5b4c11951957c6f8f642c4af61cd6b24640fec6dc7fc607ee8206a99e92410d3021ddb9a356815c3fac1026b6dec5df3124afbadb485c9ba5a3e3398a04b7ba85e58769b32a1beaf1ea27375a44095a0d1fb664ce2dd358e7fcbfb78c26a193440eb01ebfc9ed27500cd4dfc979272d1f0913cc9f66540d7e8005811109e1cf2d887c22bd8750d34016ac3c66b5ff102dacdd73f6b014e710b51e8022af9a1968ffd70157e48063fc33c97a050f7f640233bf646cc98d9524c6b92bcf3ab56f839867cc5f7f196b93bae1e27e6320742445d290f2263827498b54fec539f756afcefad4e508c098b9a7e1d8feb19955fb02ba9675585078710969d3440f5054e0f9dc3e7fe016e050eff260334f18a5d4fe391d82092319f5964f2e2eb7c1c3a5f8b13a49e282f609c317a833fb8d976d11517c571d1221a265d25af778ecf8923490c6ceeb450aecdc82e28293031d10c7d73bf85e57bf041a97360aa2c5d99cc1df82d9c4b87413eae2ef048f94b4d3554cea73d92b0f7af96e0271c691e2bb5c67add7c6caf302256adedf7ab114da0acfe870d449a3a489f781d659e8beccda7bce9f4e8618b6bd2f4132ce798cdc7a60e7e1460a7299e3c6342a579626d22733e50f526ec2fa19a22b31e8ed50f23cd1fdf94c9154ed3a7609a2f1ff981fe1d3b5c807b281e4683cc6d6315cf95b9ade8641defcb32372f1c126e398ef7a5a2dce0a8a7f68bb74560f8f71837c2c2ebbcbf7fffb42ae1896f13f7c7479a0b46a28b6f55540f89444f63de0378e3d121be09e06cc9ded1c20e65876d36aa0c65e9645644786b620e2dd2ad648ddfcbf4a7e5b1a3a4ecfe7f64667a3f0b7e2f4418588ed35a2458cffeb39b93d26f18d2ab13bdce6aee58e7b99359ec2dfd95a9c16dc00d6ef18b7933a6f8dc65ccb55667138776f7dea101070dc8796e3774df84f40ae0c8229d0d6069e5c8f39a7c299677a09d367fc7b05e3bc380ee652cdc72595f74c7b1043d0e1ffbab734648c838dfb0527d971b602bc216c9619ef0abf5ac974a1ed57f4050aa510dd9c74f508277b39d7973bb2dfccc5eeb0618db8cd74046ff337f0a7bf2c8e03e10f642c1886798d71806ab1e888d9e5ee87d0838c5655cb21c6cb83313b5a631175dff4963772cce9108188b34ac87c81c41e662ee4dd2dd7b2bc707961b1e646c4047669dcb6584f0d8d770daf5d7e7deb2e388ab20e2573d171a88108e79d820e98f26c0b84aa8b2f4aa4968dbb818ea32293237c50ba75ee485f4c22adf2f741400bdf8d6a9cc7df7ecae576221665d7358448818bb4ae4562849e949e17ac16e0be16688e156b5cf15e098c627c0056a900000000000000000000000000000000000000000000000000000000000000015b7fd3e74256b89311a8d65d3e16fc204a04270b88abb2e88fb6032920c9e01d00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000008f9113a03cf99413eed9074c13437eb5baf6ec4c000000000000000000000000000000000000000000000000002386f26fc1000000000000000000000000000000000000000000000000000000000000000005200000000000000000000000000000000000000000000000000000000000000000")
// 	toAddr := common.HexToAddress("0x2a3dd3eb832af982ec71669e178424b10dca2ede")
// 	gas := uint64(0x16e360)
// 	gasPrice := big.NewInt(100000000000)
// 	value := big.NewInt(0)
// 	v := stringToBig("0x8bd")
// 	r := stringToBig("0x3e385c86bbd28a426bc8f6df5ef51c0c7fa35d369a0aec0467b4cb8ebd194d5d")
// 	s := stringToBig("0x5ca4d4663da0f89d5b6f7b81f73ff7d77b7f5b2ad07f8509666da71c69f0541")

// 	// tx := ethtypes.NewTransaction(0, toAddr, big.NewInt(0), uint64(0x16e360), big.NewInt(0), data)
// 	txData := &ethtypes.LegacyTx{
// 		Nonce:    nonce,
// 		GasPrice: gasPrice,
// 		Gas:      gas,
// 		To:       &toAddr,
// 		Value:    value,
// 		Data:     input,
// 		V:        v,
// 		R:        r,
// 		S:        s,
// 	}
// 	tx := ethtypes.NewTx(txData)

// 	txBytes, _ := rlp.EncodeToBytes(tx)
// 	fmt.Println(txBytes)
// 	fmt.Println(len(txBytes))
// 	if len(txBytes) != 1456 {
// 		panic(len(txBytes))
// 	}

// 	signer := ethtypes.LatestSignerForChainID(big.NewInt(1101))
// 	txContentHash := signer.Hash(tx)

// 	// something = [32]byte{57, 64, 134, 223, 128, 43, 18, 142, 212, 94, 43, 113, 199, 109, 20, 172, 253, 229, 245, 25, 50, 236, 181, 41, 0, 186, 65, 70, 160, 39, 78, 80}
// 	Hash := "0x" + hex.EncodeToString(txContentHash[:])

// 	// Hash := "0x80448d99574bdd0b8df4b54dc65cdce19037464775f41daf7934fbcee5a10f7d"
// 	R := "0x3e385c86bbd28a426bc8f6df5ef51c0c7fa35d369a0aec0467b4cb8ebd194d5d"
// 	S := "0x5ca4d4663da0f89d5b6f7b81f73ff7d77b7f5b2ad07f8509666da71c69f0541"
// 	V := 0x8bd

// 	chainId := 1101

// 	Hash = Hash[2:]
// 	R = R[2:]
// 	S = S[2:]

// 	for {
// 		if len(R) == 64 {
// 			break
// 		}

// 		R = "0" + R // padding with zeroes
// 	}

// 	for {
// 		if len(S) == 64 {
// 			break
// 		}

// 		S = "0" + S // padding with zeroes
// 	}

// 	vByte := (byte)(V - (chainId*2 + 35))

// 	h, _ := hex.DecodeString(Hash)
// 	r_, _ := hex.DecodeString(R)
// 	s_, _ := hex.DecodeString(S)
// 	sig := append(r_, s_...)
// 	sig = append(sig, vByte)

// 	publicKey, err := crypto.SigToPub(h, sig)
// 	if err != nil {
// 		panic(err)
// 	}

// 	addr := crypto.PubkeyToAddress(*publicKey)
// 	if strings.ToLower(addr.Hex()) != strings.ToLower("0x8f9113a03cf99413eed9074c13437eb5baf6ec4c") {
// 		panic(fmt.Errorf("%s <-> %s", "0x8f9113a03cf99413eed9074c13437eb5baf6ec4c", addr.Hex()))
// 	}

// 	fmt.Println("Hash:", Hash)
// 	fmt.Println("Sender Address:", addr.Hex())
// }

// func stringToBig(input string) *big.Int {
// 	result := new(big.Int)
// 	result.SetString(trimHex(input), 16)
// 	return result
// }
// func trimHex(input string) string {
// 	if strings.HasPrefix(input, "0x") {
// 		return input[2:]
// 	}
// 	return input
// }
