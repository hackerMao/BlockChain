package main

import (
	"github.com/boltdb/bolt"
	"log"
)

const (
	blockChainDb  = "blockChain.db"
	blockBucket   = "blockBucket"
	lastBlockHashKey = "lastBlockHashKey"
)

type BlockChain struct {
	// 区块链数组
	//blocks []*Block
	// 使用数据库存储
	db   *bolt.DB
	tail []byte // 存储最有后一个区块hash值
}

func NewBlockChain() *BlockChain {
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
			genesisBlock := GenesisBlock()
			// 将创世块写入链中
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			// 更新最后一个区块的hash
			bucket.Put([]byte(lastBlockHashKey), genesisBlock.Hash)
			// 将lastHash写入内存中s
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

func GenesisBlock() *Block {
	return NewBlock([]byte{}, "这是区块链的第一个区块！")
}

func (self *BlockChain) AddBlock(data string) {
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
		block := NewBlock(lastHash, data)
		// 添加到区块链并更新lastHash
		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte(lastBlockHashKey), block.Hash)
		// 更新内存中区块链的最后一个hash
		self.tail = block.Hash
		return nil
	})
}
