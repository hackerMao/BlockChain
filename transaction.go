package main

// 交易模块

//noinspection GoUnresolvedReference
import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const reward = 50

// 定义交易结构
type Transaction struct {
	TXID      []byte     // 交易ID,
	TXInputs  []TXInput  //交易输入数组对整个交易的hash
	TxOutputs []TXOutput //交易输出数组
}

type TXInput struct {
	// 引用的交易ID，来自上一场交易
	TXid []byte
	// output索引值
	Index int64
	// 数字签名：由r,s拼成的[]byte
	Signature []byte
	// 公钥：由X、Y坐标拼接的字符串，在校验端重新拆分
	PublicKey []byte
}
type TXOutput struct {
	// 转账金额
	Value float64
	// 收款方的公钥hash,可通过地址逆向推理
	PubKeyHash []byte
}

func NewTXOutput(address string, value float64) *TXOutput {
	output := TXOutput{
		Value: value,
	}
	output.Lock(address)
	return &output
}

//创建公钥hash
func (self *TXOutput) Lock(address string) {
	// base58解码

	self.PubKeyHash = GetPubKeyFromAddr(address)
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
	// 矿工在挖矿时无需指定签名，故PubKeyHash可以自由填写数据，一般填写矿池的名字
	// 签名在交易完整后再填写
	input := TXInput{[]byte{}, -1, nil, []byte(data)}
	output := NewTXOutput(address, reward)
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{*output}}
	// 获取交易ID
	tx.SetHash()
	return &tx
}

func (self *Transaction) IsCoinBase() bool {
	// 判断是否是挖矿交易
	// 只有一个交易input
	// 交易ID 为空
	// 交易的index：-1
	if len(self.TXInputs) == 1 && len(self.TXInputs[0].TXid) == 0 && self.TXInputs[0].Index == -1 {
		return true
	}
	return false
}

func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	// 创建交易后需要数字签名-->需要私钥-->打开钱包“NewWallets()”
	wts := NewWallets()
	// 找到自己的钱包
	wallet := wts.WalletsMap[from]
	if wallet == nil {
		fmt.Printf("not found wallet of this address, failed to create transaction!")
		return nil
	}
	pubKey := wallet.PublicKey
	//privateKey := wallet.PrivateKey
	pubKeyHash := PublicKeyHash(pubKey)
	// 找到合理的UTXO
	utxos, resValue := bc.FindNeedUtxos(pubKeyHash, amount)
	// 与目标金额相比较，不足则返回
	if resValue < amount {
		fmt.Println("余额不足！")
		return nil
	}
	// 创建交易输入，将这些UTXO转换成input
	var inputs []TXInput
	for id, indexArray := range utxos {
		for _, i := range indexArray {
			input := TXInput{[]byte(id), int64(i), nil, pubKey}
			inputs = append(inputs, input)
		}
	}
	// 创建交易输出output
	var outputs []TXOutput
	output := NewTXOutput(to, amount)
	outputs = append(outputs, *output)
	// 找零：将剩余的转成output,转给自己
	if resValue > amount {
		output = NewTXOutput(from, resValue-amount)
		outputs = append(outputs, *output)
	}
	tx := Transaction{[]byte{}, inputs, outputs}
	tx.SetHash()
	return &tx
}
