package main

// 命令行解析库
import (
	"fmt"
	"time"
)

func (self *Cli) printChain() {
	// 新建一个迭代器来打印区块
	it := self.bc.NewIterator()
	for {
		block := it.Next()
		fmt.Printf("============================ 区块开始 ====================================================================\n")
		fmt.Printf("当前版本号：%x\n", block.Version)
		fmt.Printf("MerkalRoot：%x\n", block.MerkalRoot)
		fmt.Printf("时间戳：%s\n", time.Unix(int64(block.TimeStamp), 0).Format("2016-01-02 15:04:05"))
		fmt.Printf("当前难度值：%x\n", block.Difficulty)
		fmt.Printf("前区块哈希值：%x\n", block.PrevHash)
		fmt.Printf("当前区块哈希值：%x\n", block.Hash)
		fmt.Printf("随机数：%x\n", block.Nonce)
		fmt.Printf("数据：%s\n", block.Transactions[0].TXInputs[0].PublicKey)
		if len(block.PrevHash) == 0 {
			fmt.Println("============================ 区块遍历结束 ===============================================================")
			break
		}
	}
}

func (self *Cli) GetBalance(address string) {
	// 校验地址是否有效
	if !IsValidAddress(address) {
		fmt.Printf("地址无效：%s\n", address)
		return
	}
	// 生成公钥hash
	pubKeyHash := GetPubKeyFromAddr(address)
	fmt.Printf("=========================="+address+"==========================")
	utxos := self.bc.FindUTXOs(pubKeyHash)

	total := 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("%s的可用余额为：%f\n", address, total)
}

func (self *Cli) Send(from, to string, amount float64, miner, data string) {
	// 校验地址是否有效
	if !IsValidAddress(from) {
		fmt.Printf("地址无效：%s\n", from)
		return
	}
	if !IsValidAddress(to) {
		fmt.Printf("地址无效：%s\n", to)
		return
	}
	if !IsValidAddress(miner) {
		fmt.Printf("地址无效：%s\n", miner)
		return
	}
	// 创建挖矿交易
	coinBase := NewCoinBaseTX(miner, data)
	// 创建普通交易
	tx := NewTransaction(from, to, amount, self.bc)
	if tx == nil {
		return
	}
	// 添加区块
	self.bc.AddBlock([]*Transaction{coinBase, tx})
	fmt.Println("转账成功！明细如下：")
	fmt.Println("from: ", from)
	fmt.Println("to: ", to)
	fmt.Printf("amount: %f bitcorn\n", amount)
	fmt.Println("miner: ", miner)
	fmt.Println("data: ", data)
}

func (self *Cli) NewWallet() {
	wts := NewWallets()
	address := wts.GenerateWallets()
	fmt.Printf("钱包地址：%s\n", address)
}

func (self *Cli) ListAddresses()  {
	wts := NewWallets()
	addresses := wts.GetAllWallets()
	if len(addresses) == 0 {
		fmt.Printf("                  ")
		fmt.Printf("wallet is empty!\n")
	} else {
		for index, addr := range addresses {
			fmt.Printf("                  ")
			fmt.Printf("wallet[%d]: %s\n", index, addr)
		}
	}
}