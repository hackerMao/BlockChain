package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

const reward = 12.5

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

// 提供创建交易的方法(挖矿交易)
func NewCoinBaseTX(address string, data string) *Transaction {
	// 挖矿只有一个input
	// 无需引用交易ID
	// 无需引用👈index
	// 矿工在挖矿时无需指定签名，故sig可以自由填写数据，一般填写矿池的名字
	input := TXInput{[]byte{}, -1, data}
	output := TXOutput{reward, address}
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{output}}
	// 获取交易ID
	tx.SetHash()
	return &tx
}

// 创建挖矿交易
