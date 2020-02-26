package main

import "fmt"

type Block struct {
	PrevHash []byte
	Hash []byte
	Data []byte
}

func NewBlock(PrevBlockHash []byte, data string) *Block {
	block := Block{
		PrevHash:PrevBlockHash,
		Hash:[]byte{},
		Data:[]byte(data),
	}
	return &block
}

func main() {
	block := NewBlock([]byte{}, "可不可以给我一枚比特币？")
	fmt.Printf("前区块哈希值：%x\n", block.PrevHash)
	fmt.Printf("当前区块哈希值：%x\n", block.Hash)
	fmt.Printf("数据：%s\n", block.Data)
}
