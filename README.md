# cdk-erigon-tx-repeater

## example usage
```bash
go run ./cmd/main.go --faucet-key %YOUR_FAUCET_PRIVATE_KEY_AS_HEX_WITHOUT_0x_PREFIX% --tx-count 1000 --funding-amount 200
```
You must ensure that the faucet wallet has sufficient funds. Each account will be funded by **--funding-amount** ETH.
