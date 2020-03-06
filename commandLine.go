package main

// 命令行解析库
import "fmt"

func (self *Cli) AddBlock(data string) {
	//self.bc.AddBlock(data)
	fmt.Println("添加区块成功")
}

func (self *Cli) printChain() {
	// 新建一个迭代器来打印区块
	it := self.bc.NewIterator()
	for {
		block := it.Next()
		fmt.Printf("============================ 区块开始 ====================================================================\n")
		fmt.Printf("当前版本号：%x\n", block.Version)
		fmt.Printf("MerkalRoot：%x\n", block.MerkalRoot)
		fmt.Printf("时间戳：%x\n", block.TimeStamp)
		fmt.Printf("当前难度值：%x\n", block.Difficulty)
		fmt.Printf("前区块哈希值：%x\n", block.PrevHash)
		fmt.Printf("当前区块哈希值：%x\n", block.Hash)
		fmt.Printf("随机数：%x\n", block.Nonce)
		fmt.Printf("数据：%s\n", block.Transaction[0].TXInputs[0].Sig)
		if len(block.PrevHash) == 0 {
			fmt.Println("============================ 区块遍历结束 ===============================================================")
			break
		}
	}
}

func (self *Cli) GetBalance(address string) {
	utxos := self.bc.FindUTXOs(address)
	total := 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("余额为：%f\n", total)
}

func (self *Cli) Send(from, to string, amount float64, miner, data string) bool {
	// 创建挖矿交易
	coinBase := NewCoinBaseTX(miner, data)
	// 创建普通交易
	tx := NewTransaction(from, to, amount, self.bc)
	if tx == nil {
		return false
	}
	// 添加区块
	self.bc.AddBlock([]*Transaction{coinBase, tx})

	fmt.Println("from: ", from)
	fmt.Println("to: ", to)
	fmt.Printf("amount: %f bitcorn\n", amount)
	fmt.Println("miner: ", miner)
	fmt.Println("data: ", data)
	return true
}
