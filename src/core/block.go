package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
		"time"
)

type Block struct {
	Hash          []byte         // Hash值
	PrevBlockHash []byte         // 上一个区块的Hash值
	Timestamp     int64          // 时间戳
	Nonce         int            // 工作量证明
	Transactions  []*Transaction //交易
}

// 生成一个新的区块
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Hash:          []byte{},
		PrevBlockHash: prevBlockHash,
		Timestamp:     time.Now().Unix(),
		Nonce:         0,
		Transactions:  transactions}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

// 生成创世区块
func NewGenesisBlock(coinbase *Transaction) *Block {
	block := NewBlock([]*Transaction{coinbase}, []byte{})
	return block
}

// 区块序列化
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic("failed encode block", err)
	}
	return result.Bytes()
}

// 反序列化byte数组，生成block实例。
func Deserialize(b []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

// 计算区块里所有交易的哈希
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
