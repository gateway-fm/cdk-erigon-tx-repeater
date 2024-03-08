# cdk-erigon-tx-repeater

## example usage
```bash
go run ./cmd/main.go --faucet-key %YOUR_FAUCET_PRIVATE_KEY_AS_HEX_WITHOUT_0x_PREFIX% --blocks 1000 --funding-amount 200
```
You must ensure that the faucet wallet has sufficient funds. Each account will be funded by **--funding-amount** ETH.

Enable the following flags on your node:

zkevm.allow-free-transactions: 1
zkevm.allow-pre-eip155-transactions: 1
