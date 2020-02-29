package main

func main() {
	bc := NewBlockChain()
	bc.AddBlock("I transferred 50 bit coins to Ben")
	bc.AddBlock("Ben transferred 10 bit coins to Jerry")
	//	for index, block := range bc.blocks {
	//		fmt.Printf("======================当前区块高度：%d  ======================\n", index+1)
	//		fmt.Printf("前区块哈希值：%x\n", block.PrevHash)
	//		fmt.Printf("当前区块哈希值：%x\n", block.Hash)
	//		fmt.Printf("数据：%s\n", block.Data)
	//	}
	//	fmt.Println("======================区块结束======================")
}
