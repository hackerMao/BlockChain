package main

import (
	"crypto/sha256"
	"fmt"
)

type Block struct {
	// 前区块hash值
	PrevHash []byte
	Hash     []byte
	Data     []byte
}

func NewBlock(PrevBlockHash []byte, data string) *Block {
	block := Block{
		PrevHash: PrevBlockHash,
		Hash:     []byte{},
		Data:     []byte(data),
	}
	block.SetHash()
	return &block
}

// 生成hash
func (self *Block) SetHash() {
	// 拼装数据
	blockInfo := append(self.PrevHash, self.Data...)
	// sha256
	hash := sha256.Sum256(blockInfo)
	self.Hash = hash[:]
}

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

func main() {
	bc := NewBlockChain()
	bc.AddBlock("我向MM转了50枚比特币")
	bc.AddBlock("我向MM转了100枚比特币")
	for index, block := range bc.blocks {
		fmt.Printf("======================当前区块高度：%d  ======================\n", index+1)
		fmt.Printf("前区块哈希值：%x\n", block.PrevHash)
		fmt.Printf("当前区块哈希值：%x\n", block.Hash)
		fmt.Printf("数据：%s\n", block.Data)
	}
	fmt.Println("======================区块结束======================")
}
