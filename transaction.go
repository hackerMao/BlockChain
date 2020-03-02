package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

// 定义交易结构
type Transaction struct {
	TXID     []byte     // 交易ID,对整个交易的hash
	TXInput  []TXInput  //交易输入数组
	TxOutput []TXOutput //交易输出数组
}

type TXInput struct {
	// 引用的交易ID
	TXid []byte
	// output索引值
	Index int64
	// 解锁脚本，使用地址模拟
	Sig string
}
type TXOutput struct {
	// 转账金额
	value float64
	// 锁定脚本，使用地址模拟
	PubKeyHash string
}

func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	// 新建一个编码器
	encoder := gob.NewEncoder(&buffer)
	// 对交易编码
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}

// 提供创建交易的方法
// 创建挖矿交易
