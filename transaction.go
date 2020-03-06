package main
// 交易模块

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const reward = 12.5

// 定义交易结构
type Transaction struct {
	TXID     []byte     // 交易ID,对整个交易的hash
	TXInputs  []TXInput  //交易输入数组
	TxOutputs []TXOutput //交易输出数组
}

type TXInput struct {
	// 引用的交易ID，来自上一场交易
	TXid []byte
	// output索引值
	Index int64
	// 解锁脚本，使用地址模拟
	Sig string
}
type TXOutput struct {
	// 转账金额
	Value float64
	// 锁定脚本，使用地址模拟
	PubKeyHash string
}

//设置交易ID
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
	// 挖矿交易只有一个input
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

func (self *Transaction) IsCoinBase() bool {
	// 判断是否是挖矿交易
	// 只有一个交易input
	if len(self.TXInputs) == 1 {
		// 交易ID 为空
		// 交易的index：-1
		input := self.TXInputs[0]
		if bytes.Equal(input.TXid, []byte{}) && input.Index == -1 {}
		return true
	}
	return false
}

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	// 找到合理的UTXO

	utxos, resValue := bc.FindNeedUtxos(from, amount)
	// 与目标金额相比较，不足则返回
	if resValue < amount {
		fmt.Println("余额不足！")
		return nil
	}
	// 创建交易输入，将这些UTXO转换成input
	var inputs []TXInput
	for id, indexArray := range utxos {
		for _, i := range indexArray {
			input := TXInput{[]byte(id), int64(i), from}
			inputs = append(inputs, input)
		}
	}
	// 创建交易输出output
	var outputs []TXOutput
	outputs = append(outputs, TXOutput{amount, to})
	// 找零：将剩余的转成output,转给自己
	if resValue > amount {
		outputs = append(outputs, TXOutput{resValue-amount, from})
	}
	tx := Transaction{[]byte{}, inputs, outputs}
	tx.SetHash()
	fmt.Println("tx:", tx)
	return &tx
}