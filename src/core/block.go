package core

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Data          []byte // 存放数据
	Hash          []byte // Hash值
	PrevBlockHash []byte // 上一个区块的Hash值
	Timestamp     int64  // 时间戳
}

// 生成一个新的区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{[]byte(data), []byte{}, prevBlockHash, time.Now().Unix()}
	block.SetHash()
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
