package main

//区块链
import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

const (
	blockChainDb     = "blockChain.db"
	blockBucket      = "blockBucket"
	lastBlockHashKey = "lastBlockHashKey"
)

type BlockChain struct {
	// 区块链数组
	//blocks []*Block
	// 使用数据库存储
	db   *bolt.DB
	tail []byte // 存储最有后一个区块hash值
}

func NewBlockChain(address string) *BlockChain {
	var lastHash []byte
	db, err := bolt.Open(blockChainDb, 0600, nil)
	if err != nil {
		log.Panic("连接数据库失败")
	}
	// 更新或创建数据库
	db.Update(func(tx *bolt.Tx) error {
		// 找到bucket抽屉，如果没有则创建一个
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("创建bucket失败")
			}
			// 创建第一个区块--创世块
			genesisBlock := GenesisBlock(address)
			// 将创世块写入链中
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			// 更新最后一个区块的hash
			bucket.Put([]byte(lastBlockHashKey), genesisBlock.Hash)
			// 将lastHash写入内存中
			lastHash = genesisBlock.Hash

			// 测试
			//blockBytes := bucket.Get(genesisBlock.Hash)
			//block := Deserialize(blockBytes)
			//fmt.Println(block)
		} else {
			lastHash = bucket.Get([]byte(lastBlockHashKey))
		}
		return nil
	})
	return &BlockChain{
		db:   db,
		tail: lastHash,
	}
}

func GenesisBlock(address string) *Block {
	coinBase := NewCoinBaseTX(address, "创世块天下第一")
	return NewBlock([]byte{}, []*Transaction{coinBase})
}

func (self *BlockChain) AddBlock(txs []*Transaction) {
	// 连接数据库
	db := self.db
	lastHash := self.tail
	// 更新
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("数据库bucket异常，请检查")
		}
		// 创建区块
		block := NewBlock(lastHash, txs)
		// 添加到区块链并更新lastHash
		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte(lastBlockHashKey), block.Hash)
		// 更新内存中区块链的最后一个hash
		self.tail = block.Hash
		return nil
	})
}
// 查找所有可用的UTXO
func (self *BlockChain) FindUTXOs(address string) []TXOutput {
	var UTXO []TXOutput
	//定义一个map来保存以消费掉的output
	spendOutputs := make(map[string][]int64)
	// 创建区块链迭代器
	it := self.NewIterator()
	// 遍历所有区块
	for {
		block := it.Next()
		// 遍历区块的交易
		for _, tx := range block.Transaction {
			// 遍历交易中的所有output
			if !tx.IsCoinBase() {
				// 遍历input,找到自己花费过的utxo的集合
				for _, input := range tx.TXInputs {
					if input.Sig == address {
						// 记录交易ID和该场交易下属于我的花费
						arrayIndex := spendOutputs[string(tx.TXID)]
						arrayIndex = append(arrayIndex, input.Index)
					}
				}
			}
		OUTPUT:
			for i, output := range tx.TxOutputs {
				//fmt.Println(output)
				// 过滤掉已花费的output
				if spendOutputs[string(tx.TXID)] != nil {
					// 通过交易ID 找到该次交易花费过的utxo索引值
					for _, j := range spendOutputs[string(tx.TXID)] {
						if int64(i) == j {
							continue OUTPUT
						}
					}
				}
				// 添加满足条件的output
				if output.PubKeyHash == address {
					UTXO = append(UTXO, output)
				}
			}
		}
		if len(block.PrevHash) == 0 {
			fmt.Printf("区块遍历结束\n")
			break
		}
	}
	return UTXO
}

func (self *BlockChain) FindNeedUtxos(from string, amount float64) (map[string][]uint64, float64) {
	utxos := make(map[string][]uint64)
	var calc float64
	//定义一个map来保存以消费掉的output
	spendOutputs := make(map[string][]int64)
	// 创建区块链迭代器
	it := self.NewIterator()
	// 遍历所有区块
	I:
		for {
		block := it.Next()
		// 遍历区块的交易
		for _, tx := range block.Transaction {
			// 遍历交易中的所有output
		OUTPUT:
			for i, output := range tx.TxOutputs {
				//fmt.Println(output)
				// 过滤掉已花费的output
				if spendOutputs[string(tx.TXID)] != nil {
					// 通过交易ID 找到该次交易花费过的utxo索引值
					for _, j := range spendOutputs[string(tx.TXID)] {
						if int64(i) == j {
							continue OUTPUT
						}
					}
				}
				// 添加满足条件的output
				if output.PubKeyHash == from {
					// 找到自己需要的utxo
					if calc < amount {
						utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], uint64(i))
						calc += output.Value

						if calc >= amount {
							break I
						}
					}
				}
			}
			if !tx.IsCoinBase() {
				// 遍历input,找到自己花费过的utxo的集合
				for _, input := range tx.TXInputs {
					if input.Sig == from {
						// 记录交易ID和该场交易下属于我的花费
						arrayIndex := spendOutputs[string(tx.TXID)]
						arrayIndex = append(arrayIndex, input.Index)
					}
				}
			}
		}
		if len(block.PrevHash) == 0 {
			fmt.Printf("区块遍历结束\n")
			break I
		}
	}
	return utxos, calc
}