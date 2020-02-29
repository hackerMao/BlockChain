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
	// 获取区块链中最后一个区块，取该区块的哈希值作为上一个哈希值
	//lastlock := self.blocks[len(self.blocks)-1]
	//// 新建一个区块
	//block := NewBlock(lastlock.Hash, data)
	//// 添加区块
	//self.blocks = append(self.blocks, block)
}
