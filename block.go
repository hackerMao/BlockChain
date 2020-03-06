package main
// 区块
import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	// 版本号
	Version uint64
	// 前区块hash值
	PrevHash []byte
	// Merkal根（梅卡尔根），交易时产生的hash值:对交易进行二叉树hash
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
	//Data []byte
	Transaction []*Transaction
}

func (self *Block) Serialize() []byte {
	//新建一个编码器
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	//使用编码器对block进行编码
	err := encoder.Encode(self)
	if err != nil {
		log.Panic("编码失败:",err)
	}
	return buffer.Bytes()
}

// 模拟梅克尔根，只是对交易的数据作简单的拼接，而不做二叉树处理
func (self *Block) MakeMerkalRoot() []byte {
	//TODO
	return []byte{}
}

func Deserialize(b []byte) Block {
	// 新建一个解码器
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(b))
	// 使用解码器对bytes进行解码
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("解码失败")
	}
	return block
}

func NewBlock(PrevBlockHash []byte, txs []*Transaction) *Block {
	block := Block{
		Version:     00,
		PrevHash:    PrevBlockHash,
		MerkalRoot:  []byte{},
		TimeStamp:   uint64(time.Now().Unix()),
		Difficulty:  0,
		Nonce:       0,
		Hash:        []byte{},
		Transaction: txs,
	}
	block.MerkalRoot = block.MakeMerkalRoot()
	// 创建一个工作量证明
	pow := NewProofWork(&block)
	// 查找随机数、不停进行hash运算。
	hash, nonce := pow.Run()
	// 更新区块信息
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

/*
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
*/

func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}
