package main
// 区块链迭代器

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockChainIterator struct {
	db                 *bolt.DB
	CurrentHashPointer []byte
}

func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{
		bc.db,
		bc.tail}
}


func (self *BlockChainIterator) Next() *Block {
	var block Block
	self.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("未知异常，要操作的bucket不存在")
		}
		// 读取当前游标指向的区块
		blockBytes := bucket.Get(self.CurrentHashPointer)
		// 反序列化
		block = Deserialize(blockBytes)
		// 游标往左移
		self.CurrentHashPointer = block.PrevHash
		return nil
	})
	return &block
}