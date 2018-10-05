package core

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

const (
	DB_FILE               = "../db/BlockChain.db"
	BLOCKS_BUCKET         = "blocks"
	LAST_HASH             = "lastHash"
	GENESIS_COINBASE_DATA = "The Times 03/October/2018 Chancellor on brink of second bailout for banks"
)

type BlockChain struct {
	lastHash []byte
	db       *bolt.DB
}

type BlockChainIterator struct {
	currentHash []byte
	DB          *bolt.DB
}

// 得到区块链
func GetBlockChain(address string) *BlockChain {
	if dbExists() == false {
		fmt.Println("No existing block chain found. Create one first.")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(DB_FILE, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKS_BUCKET))
		tip = b.Get([]byte(LAST_HASH))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bc := BlockChain{tip, db}
	return &bc
}

func dbExists() bool {
	if _, err := os.Stat(DB_FILE); os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateBlockchain 创建一个新的区块链数据库
// address 用来接收挖出创世块的奖励
func CreateBlockchain(address string) *BlockChain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(DB_FILE, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		// 挖矿奖励交易
		cbtx := NewCoinbaseTX(address, GENESIS_COINBASE_DATA)
		// 生成创世区块
		genesis := NewGenesisBlock(cbtx)
		b, err := tx.CreateBucket([]byte(BLOCKS_BUCKET))
		// 将区块放入数据库
		err = b.Put(genesis.Hash, genesis.Serialize())
		err = b.Put([]byte(LAST_HASH), genesis.Hash)
		tip = genesis.Hash
		return err
	})
	if err != nil {
		log.Panic(err)
	}
	bc := BlockChain{tip, db}
	return &bc
}

// 挖矿
func (bc *BlockChain) MineBlock(transactions []*Transaction) {
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
	newBlock := NewBlock(transactions, lastHash)
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
		log.Panic("failed to add new block:", err)
	}
}

// 关闭数据库连接
func (bc *BlockChain) Close() {
	bc.db.Close()
}

// 生成Iterator实例
func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.lastHash, bc.db}
}

// 遍历所有区块
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

// 找到未花费输出的交易
func (bc *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()
	for {
		block := bci.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			for outIdx, out := range tx.Vout {
				// 如果交易输出被花费了
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				// 如果该交易输出可以被解锁，即可被花费
				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
}

func (bc *BlockChain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(address)
	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

// 从 address 中找到至少 amount 的 UTXO
func (bc *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0
Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOutputs
}
