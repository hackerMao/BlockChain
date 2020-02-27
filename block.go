package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"time"
)

type Block struct {
	// 版本号
	Version uint64
	// 前区块hash值
	PrevHash []byte
	// Merkal根（梅卡尔根），交易时产生呢的hash值
	MerkalRoot []byte
	// 时间戳
	TimeStamp uint64
	// 难度值
	Difficulty uint64
	// 随机数
	Nonce uint64
	// 当前哈希值，正常情况下比特币中没有当前hash值
	Hash []byte
	// 数据
	Data []byte
}

func NewBlock(PrevBlockHash []byte, data string) *Block {
	block := Block{
		Version:    00,
		PrevHash:   PrevBlockHash,
		MerkalRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce:      0,
		Hash:       []byte{},
		Data:       []byte(data),
	}
	block.SetHash()
	return &block
}

// 生成hash
func (self *Block) SetHash() {
	var blockInfo []byte
	// 拼装数据
	blockInfo = append(blockInfo, Uint64ToByte(self.Version)...)
	blockInfo = append(blockInfo, self.PrevHash...)
	blockInfo = append(blockInfo, self.MerkalRoot...)
	blockInfo = append(blockInfo, Uint64ToByte(self.TimeStamp)...)
	blockInfo = append(blockInfo, Uint64ToByte(self.Difficulty)...)
	blockInfo = append(blockInfo, Uint64ToByte(self.Nonce)...)
	blockInfo = append(blockInfo, self.Data...)
	// sha256
	hash := sha256.Sum256(blockInfo)
	self.Hash = hash[:]
}

func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}
