package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Panic("连接数据库失败")
	}
	// 更新或创建数据库
	db.Update(func(tx *bolt.Tx) error {
		// 找到bucket抽屉，如果没有则创建一个
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte("b1"))
			if err != nil {
				log.Panic("创建bucket(b1)失败")
			}
		}
		bucket.Put([]byte("name"), []byte("hacker_murray"))
		bucket.Put([]byte("age"), []byte("24"))
		return nil
	})
	// 读取数据
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil {
			log.Panic("数据库不存在")
		}
		name := bucket.Get([]byte("name"))
		age := bucket.Get([]byte("age"))
		fmt.Printf("name:%s, age:%s\n", string(name), string(age))
		return nil
	})
}
