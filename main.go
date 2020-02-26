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
	return  &BlockChain{
		blocks:[]*Block{genesisBlock},
	}
}

func GenesisBlock() *Block {
	return NewBlock([]byte{}, "这是区块链的第一个区块！")
}

func main() {
	bc := NewBlockChain()
	for index, block := range bc.blocks {
		fmt.Printf("======================当前区块高度：%d  ======================\n", index+1)
		fmt.Printf("前区块哈希值：%x\n", block.PrevHash)
		fmt.Printf("当前区块哈希值：%x\n", block.Hash)
		fmt.Printf("数据：%s\n", block.Data)
	}
	fmt.Println("======================区块结束======================")
}
