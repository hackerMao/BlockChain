package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "wallet.dat"

type Wallets struct {
	// map[地址]钱包
	WalletsMap map[string]*Wallet
}

func NewWallets() *Wallets {
	var wts Wallets
	wts.WalletsMap = make(map[string]*Wallet)
	wts.loadFile()
	return &wts
}

// 创建钱包组
func (wts *Wallets) GenerateWallet() string {
	wt := NewWallet()
	address := wt.GenerateAddr()
	wts.WalletsMap[address] = wt
	wts.saveToFile()
	return address
}

// 将钱包组信息保存到文件
func (wts *Wallets) saveToFile() {
	var buffer bytes.Buffer
	gob.Register(elliptic.P256())
	// 新建一个编码器
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(wts)
	if err != nil {
		log.Panic(err)
	}
	ioutil.WriteFile(walletFile, buffer.Bytes(), 0600)
}

// 读取文件，把所有的wallet读出
func (wts *Wallets) loadFile()  {
	// 在读取之前要先确认文件是否存在，如果不存在则直接退出
	_, err := os.Stat(walletFile)
	if os.IsNotExist(err) {
		return
	}
	// 读取内容
	content, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	// 解码
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	var wtsLocal Wallets
	err = decoder.Decode(&wtsLocal)
	if err != nil {
		log.Panic(err)
	}
	// 将读出来的钱包赋值给当前对象wts的钱包
	wts.WalletsMap = wtsLocal.WalletsMap
}

func (wts *Wallets) GetAllWallets()  []string {
	var addresses []string
	for address := range wts.WalletsMap {
		addresses = append(addresses, address)
	}
	return addresses
}