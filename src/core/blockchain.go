package core

type BlockChain struct {
	Blocks []*Block
}

// 生成一个新的区块链
func NewBlockChain() *BlockChain {
	bc := &BlockChain{[]*Block{NewGenesisBlock()}}
	return bc
}

// 加入一个新的区块到区块链中
func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}
