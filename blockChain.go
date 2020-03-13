package main

//区块链
import (
	"bytes"
	"crypto/ecdsa"
	"errors"
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

func (bc *BlockChain) AddBlock(txs []*Transaction) {
	// 连接数据库
	db := bc.db
	lastHash := bc.tail
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
		bc.tail = block.Hash
		return nil
	})
}

// 查找所有可用的UTXO
func (bc *BlockChain) FindUTXOs(pubKeyHash []byte) []TXOutput {
	var UTXO []TXOutput

	txs := bc.FindUTXOTransactions(pubKeyHash)
	for _, tx := range txs {
		for _, output := range tx.TxOutputs {
			if bytes.Equal(pubKeyHash, output.PubKeyHash) {
				UTXO = append(UTXO, output)
			}
		}
	}
	return UTXO
}

func (bc *BlockChain) FindNeedUtxos(senderPubHash []byte, amount float64) (map[string][]uint64, float64) {
	//找到够用的utxo集合
	utxos := make(map[string][]uint64)
	var calc float64

	// 找到所有可用的交易
	txs := bc.FindUTXOTransactions(senderPubHash)
	// 遍历交易
	for _, tx := range txs {
		// 遍历output
		for index, output := range tx.TxOutputs {
			if bytes.Equal(senderPubHash, output.PubKeyHash) {
				if calc < amount {
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], uint64(index))
					calc += output.Value
					// 找到满足的金额返回
					if calc >= amount {
						return utxos, calc
					}
				}
			}
		}
	}
	return utxos, calc
}

func (bc *BlockChain) FindUTXOTransactions(senderPubHash []byte) []*Transaction {
	var txs []*Transaction // 存储所有包含utxo的集合
	// 定义一个map来保存消费过得output，key为交易ID，value为索引值
	spentOutputs := make(map[string][]int64)

	// 创建区块链迭代器
	it := bc.NewIterator()
	for {
		// 遍历区块
		block := it.Next()
		// 遍历交易
		for _, tx := range block.Transactions {
			// 如果当前交易是挖矿交易则不做遍历
			if !tx.IsCoinBase() {
				for _, input := range tx.TXInputs {
					PubKeyHash := PublicKeyHash(input.PublicKey)
					if bytes.Equal(PubKeyHash, senderPubHash) {
						spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index)
					}
				}
			} else {
			}
		OUTPUT:
			// 遍历交易中的output，找到和address相关的，添加之前检查是否已经被消费过
			for index, output := range tx.TxOutputs {
				if spentOutputs[string(tx.TXID)] != nil {
					for _, output_index := range spentOutputs[string(tx.TXID)] {
						if output_index == int64(index) {
							continue OUTPUT
						}
					}
				}
				if bytes.Equal(output.PubKeyHash, senderPubHash) {
					txs = append(txs, tx)
				}
			}
		}
		if len(block.PrevHash) == 0 {
			fmt.Println("区块遍历结束")
			break
		}
	}
	return txs
}

func (bc *BlockChain) FindTransactionByTXid(id []byte) (Transaction, error) {
	it := bc.NewIterator()

	for {
		block := it.Next()
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.TXID, id) {
				return *tx, nil
			}
		}
		if len(block.PrevHash) == 0 {
			fmt.Printf("区块链遍历结束！\n")
			break
		}
	}
	return Transaction{}, errors.New("无效的交易ID！")
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey *ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)
	// 找到所有引用的交易
	// 遍历inputs,找到目标交易
	// 添加到prevTXs
	for _, input := range tx.TXInputs {
		tx, err := bc.FindTransactionByTXid(input.TXid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[string(input.TXid)] = tx
	}
	//签名
	tx.Sign(privateKey, prevTXs)
}