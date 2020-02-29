package main

import "fmt"

func main() {
	bc := NewBlockChain()
	bc.AddBlock("I transferred 50 bit coins to Ben")
	bc.AddBlock("Ben transferred 10 bit coins to Jerry")

	// 创建一个迭代器
	it := bc.NewIterator()
	for {
		block := it.Next()
		fmt.Printf("============================ 区块开始 ====================================================================\n")
		fmt.Printf("前区块哈希值：%x\n", block.PrevHash)
		fmt.Printf("当前区块哈希值：%x\n", block.Hash)
		fmt.Printf("数据：%s\n", block.Data)
		if len(block.PrevHash) == 0 {
			fmt.Println("============================ 区块遍历结束 ===============================================================")
			break
		}
	}
}
