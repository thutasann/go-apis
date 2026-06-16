# Blockchain in Golang

## Architecture

```mermaid
flowchart LR
	Client[Client / CLI] --> Node[Node]
	Client -->|submit tx| TxPool[Transactions / TxPool]
	Node --> Blockchain[Blockchain]
	TxPool -->|included in block| Blockchain
	Blockchain --> Storage[(Storage)]
	Blockchain --> Network[Network]
	Blockchain --> Crypto[Crypto]
	Blockchain --> HeadHash["Head Hash (tip)"]
	Storage --> HeadHash
```

## Code flow (short)

- **Start:** `main.go` boots the node/server (`network/server.go`).
- **Submit tx:** client → TxPool (`core/transaction.go`).
- **Create block:** node collects txs → `core/block.go`.
- **Validate:** `core/validator.go` checks block/tx rules.
- **Hashing:** `core/hasher.go` computes header/hash.
- **Add to chain:** `core/blockchain.go` appends block and updates head/tip.
- **Persist:** `core/storage.go` stores headers/blocks and head hash.
- **Network:** peers exchange headers/blocks (`network/transport.go`).

## Component Hierarchy

```mermaid
graph TD
	App[Application]
	App --> Node[Node/Server]
	Node --> Core[Core]
	Core --> Block[Block]
	Core --> Blockchain[Blockchain]
	Core --> Tx[Transaction]
	Core --> Validator[Validator]
	Core --> Hasher[Hasher]
	Core --> Storage[Storage]
	App --> Network[Network]
	App --> Crypto[Crypto]
```
