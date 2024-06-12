# cdk-erigon-tx-repeater

## Setup
Let's **source** stands for the network where we are going to read txs from and **target** is our sequencer

Ensure both network are using the same chain-id.

Ensure that that you have the private-key of a wallet with enough funds on target network.

Run go run ./cmd/main.go --help (TxRepeater) to see the list of flags.

You must specify almost all of them.
- --source-datastream-endpoint
- --source-rpc-endpoint
- --target-rpc-endpoint (the default is http://localhost:8467)
- --faucet-key (private key of the faucet w/o 0x prefix)
- --blocks (how many from source network will be repeated on target network)
- --funding-amount (how many eth will be send to each wallet, default is 200)
- --chain-id (chain id of both networks)

**Notes:**

Before each execution of the TxRepeater you must initialize the sequencer from scratch by deleting the data folder

Ensure you have allowed pre-EIP155 transactions and free transactions using correcponding flags of the node

## Example usage
```bash
go run ./cmd/main.go --faucet-key %YOUR_FAUCET_PRIVATE_KEY_AS_HEX_WITHOUT_0x_PREFIX% --blocks 1000 --funding-amount 200
```
You must ensure that the faucet wallet has sufficient funds. Each account will be funded by **--funding-amount** ETH.

Enable the following flags on your node:

zkevm.allow-free-transactions: 1
zkevm.allow-pre-eip155-transactions: 1
