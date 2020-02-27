package main

import "crypto/sha256"

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