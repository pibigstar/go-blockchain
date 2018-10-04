package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

type Block struct {
	Data          []byte // 存放数据
	Hash          []byte // Hash值
	PrevBlockHash []byte // 上一个区块的Hash值
	Timestamp     int64  // 时间戳
	Nonce         int    // 工作量证明
}

// 生成一个新的区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{[]byte(data), []byte{}, prevBlockHash, time.Now().Unix(), 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

// 设置Hash值
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{timestamp, b.Data, b.Hash}, []byte{})
	newHash := sha256.Sum224(headers)
	b.Hash = newHash[:]
}

// 生成创世区块
func NewGenesisBlock() *Block {
	block := NewBlock("Genesis Block", []byte{})
	return block
}

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
