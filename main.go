package main

import (
	"crypto/sha256"
	"fmt"
)

type Block struct {
	// 前区块hash值
	PrevHash []byte // TODO：先留空，后边计算
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

func main() {
	block := NewBlock([]byte{}, "可不可以给我一枚比特币？")
	fmt.Printf("前区块哈希值：%x\n", block.PrevHash)
	fmt.Printf("当前区块哈希值：%x\n", block.Hash)
	fmt.Printf("数据：%s\n", block.Data)
	//fmt.Println(len("0935d481a8c7db864698443edcf84268f784a93cba5c63b7587e7c6cf60668c7"))
}
