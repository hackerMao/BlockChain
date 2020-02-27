package main

type BlockChain struct {
	// 区块链数组
	blocks []*Block
}

func NewBlockChain() *BlockChain {
	genesisBlock := GenesisBlock()
	return &BlockChain{
		blocks: []*Block{genesisBlock},
	}
}

func GenesisBlock() *Block {
	return NewBlock([]byte{}, "这是区块链的第一个区块！")
}

func (self *BlockChain) AddBlock(data string) {
	// 获取区块链中最后一个区块，取该区块的哈希值作为上一个哈希值
	lastlock := self.blocks[len(self.blocks)-1]
	// 新建一个区块
	block := NewBlock(lastlock.Hash, data)
	// 添加区块
	self.blocks = append(self.blocks, block)
}
