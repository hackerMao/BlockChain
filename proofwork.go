package main
// 工作量证明/挖矿

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	block *Block
	// 目标哈希值
	target *big.Int
}

func NewProofWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000000"
	tempInt := big.Int{}
	tempInt.SetString(targetStr, 16)
	pow.target = &tempInt
	return &pow
}

func (pow *ProofOfWork) Run() ([]byte, uint64) {
	// 生成随机数
	var hash [32]byte
	var nonce uint64
	block := pow.block
	// 生成hash值
	fmt.Println("==================开始挖矿==================")
	for {
		temp := [][]byte{
			Uint64ToByte(block.Version),
			block.PrevHash,
			Uint64ToByte(block.Difficulty),
			block.MerkalRoot,
			Uint64ToByte(nonce),
			Uint64ToByte(block.TimeStamp),
		}
		//fmt.Println(nonce)
		blockInfo := bytes.Join(temp, []byte{})
		hash := sha256.Sum256(blockInfo)
		// 与目标值进行比较
		tempInt := big.Int{}
		tempInt.SetBytes(hash[:])
		if tempInt.Cmp(pow.target) == -1 {
			fmt.Printf("挖矿成功！hash value: %x,  nonce:%d \n", hash, nonce)
			return hash[:], nonce
		} else {
			// 未找到，随机数加1，继续hash
			nonce++
		}
	}
	return hash[:], nonce
}
