package core

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Transaction 交易記錄
type Transaction struct {
	Sender    string  `json:"sender,omitempty"`
	Recipient string  `json:"recipient,omitempty"`
	Amount    float64 `json:"amount,omitempty"`
}

// Block 區塊
type Block struct {
	Index        int            `json:"index,omitempty"`
	Timestamp    time.Time      `json:"timestamp,omitempty"`
	Transactions []*Transaction `json:"transactions,omitempty"`
	Proof        int            `json:"proof,omitempty"`
	PreviousHash string         `json:"previous_hash,omitempty"`
}

// Blockchain 區塊鍊
type Blockchain struct {
	Chain               []*Block        `json:"chain,omitempty"`
	Nodes               map[string]bool `json:"nodes,omitempty"`
	CurrentTransactions []*Transaction  `json:"current_transactions,omitempty"`
}

// NewBlockchain create Blockchain
func NewBlockchain() *Blockchain {
	bc := &Blockchain{}

	bc.Nodes = make(map[string]bool)

	// Create the genesis block
	bc.NewBlock()

	return bc
}

// NewBlock 建立新的區塊
func (bc *Blockchain) NewBlock() error {
	var previousHash string
	var err error

	if len(bc.Chain) == 0 {
		previousHash = "1"
	} else {
		previousHash, err = bc.Chain[len(bc.Chain)-1].hash()
		if err != nil {
			return errors.Wrap(err, "bc.Chain[len(bc.Chain)-1].hash")
		}
	}

	b := &Block{
		Index:        len(bc.Chain) + 1,
		Timestamp:    time.Now(),
		Transactions: bc.CurrentTransactions,
		PreviousHash: previousHash,
	}

	b.proofWork()

	bc.CurrentTransactions = []*Transaction{}

	bc.Chain = append(bc.Chain, b)

	return nil
}

//NewTransaction 建立交易
func (bc *Blockchain) NewTransaction(sender, recipient string, amount float64) {
	t := &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}

	bc.CurrentTransactions = append(bc.CurrentTransactions, t)
}

//RegisterNode 註冊節點，節點用來更新區塊鍊
func (bc *Blockchain) RegisterNode(address string) {
	node := address
	bc.Nodes[node] = true
}

//validChain 驗證區塊鍊
func (bc *Blockchain) validChain() (bool, error) {
	var previousHash string

	for i, b := range bc.Chain {
		ok, err := b.validProof()
		if err != nil {
			return false, errors.Wrap(err, "b.validProof")
		}
		if !ok {
			return false, nil
		}

		if i == 0 {
			previousHash, err = b.hash()
			if err != nil {
				return false, errors.Wrap(err, "b.hash")
			}
			continue
		}

		if previousHash != b.PreviousHash {
			return false, nil
		}

		previousHash, err = b.hash()
		if err != nil {
			return false, errors.Wrap(err, "b.hash")
		}
	}

	return true, nil
}

//ResolveConflicts 與其他節點比對，確認目前正確的區塊鍊
func (bc *Blockchain) ResolveConflicts(newbc *Blockchain) (bool, error) {
	ok, err := newbc.validChain()
	if err != nil {
		return false, errors.Wrap(err, "newbc.validChain")
	}
	if ok && len(newbc.Chain) > len(bc.Chain) {
		bc.Chain = make([]*Block, len(newbc.Chain))
		copy(bc.Chain, newbc.Chain)

		return true, nil
	}

	return false, nil
}

func (bc *Blockchain) String() string {
	return fmt.Sprintf("chain:%+v\nnodes:%+v\ncurrentTransactions:%+v", bc.Chain, bc.Nodes, bc.CurrentTransactions)
}

// hash sha256
func (b *Block) hash() (string, error) {
	s, err := json.Marshal(b)
	if err != nil {
		return "", errors.Wrap(err, "json.Marshal")
	}
	sum := sha256.Sum256(s)

	return fmt.Sprintf("%x", sum), nil
}

// 超簡單驗證
func (b *Block) validProof() (bool, error) {
	h, err := b.hash()
	if err != nil {
		return false, errors.Wrap(err, "b.hash")
	}

	return h[:1] == "0", nil
}

//proofWork 挖礦工作
func (b *Block) proofWork() error {
	b.Proof = 0
	for {
		ok, err := b.validProof()
		if err != nil {
			return errors.Wrap(err, "b.validProof")
		}
		if ok {
			break
		}
		b.Proof++
	}
	return nil
}

func (b *Block) String() string {
	return fmt.Sprintf("{\n#%d %s\ntransactions:%+v\nproof:%d\nprevious Hash:%s\n}", b.Index,
		b.Timestamp,
		b.Transactions,
		b.Proof,
		b.PreviousHash,
	)
}

func (t *Transaction) String() string {
	return fmt.Sprintf("{s:%s, r:%s, $%.2f}", t.Sender, t.Recipient, t.Amount)
}
