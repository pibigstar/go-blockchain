package core

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

const (
	DB_FILE       = "../db/blockchain.db"
	BLOCKS_BUCKET = "blocks"
	LAST_HASH     = "lastHash"
)

type BlockChain struct {
	lastHash []byte
	db       *bolt.DB
}

type BlockChainIterator struct {
	currentHash []byte
	DB          *bolt.DB
}

// 生成一个新的区块链
func NewBlockChain() *BlockChain {
	var lastHash []byte
	db, err := bolt.Open(DB_FILE, 0600, nil)
	if err != nil {
		log.Panic("failed to open the db", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		//
		b := tx.Bucket([]byte(BLOCKS_BUCKET))
		if b == nil {
			fmt.Println("No existing BlockChain. Creating a new...")
			b, err := tx.CreateBucket([]byte(BLOCKS_BUCKET))
			if err != nil {
				log.Panic("failed to create bucket")
			}
			genesis := NewGenesisBlock()
			// Key - Value
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panicf("failed put hash: %x", genesis.Hash)
			}
			err = b.Put([]byte(LAST_HASH), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			lastHash = genesis.Hash
		} else {
			lastHash = b.Get([]byte(LAST_HASH))
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	bc := &BlockChain{lastHash, db}
	return bc
}

// 加入一个新的区块到区块链中
func (bc *BlockChain) AddBlock(data string) {
	var (
		lastHash []byte
		err      error
	)
	// 获取最后区块的hash值
	err = bc.db.View(func(tx1 *bolt.Tx) error {
		b := tx1.Bucket([]byte(BLOCKS_BUCKET))
		lastHash = b.Get([]byte(LAST_HASH))
		return nil
	})
	// 计算新的区块
	newBlock := NewBlock(data, lastHash)
	bc.lastHash = newBlock.Hash
	// 将新的区块加入到数据库区块链中
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKS_BUCKET))
		err = b.Put(newBlock.Hash, newBlock.Serialize())
		// 将新的区块置为最后一个区块
		err = b.Put([]byte(LAST_HASH), newBlock.Hash)
		return err
	})
	if err != nil {
		log.Panic("failed to add new block:",err)
	}
}

func (bc *BlockChain) Close() {
	bc.db.Close()
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKS_BUCKET))
		lastHash = b.Get([]byte(LAST_HASH))
		return nil
	})
	if err != nil {
		log.Println("failed view the db:", err)
	}
	return &BlockChainIterator{lastHash, bc.db}
}

func (bci *BlockChainIterator) Next() *Block {
	var byteBlock []byte
	err := bci.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKS_BUCKET))
		byteBlock = b.Get(bci.currentHash)
		return nil
	})
	if err != nil {
		log.Println("failed view the db:", err)
	}
	block := Deserialize(byteBlock)
	bci.currentHash = block.PrevBlockHash
	return block
}
