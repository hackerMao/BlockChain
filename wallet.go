package main

import (
	"BlockChain/base58"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

type Wallet struct {
	// 私钥
	PrivateKey *ecdsa.PrivateKey
	// 公钥,PublicKey不存储原始的公钥，而是存储X、Y坐标拼接的字符串，参考r、s的传递
	PublicKey []byte
}

// 创建钱包
func NewWallet() *Wallet {
	// 创建一个椭圆曲线
	Carve := elliptic.P256()
	// 通过ecdsa生成密钥对
	PrivateKey, err := ecdsa.GenerateKey(Carve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	// 获取公钥
	PublicKey := PrivateKey.PublicKey
	// X、Y坐标拼接的字符串
	publicKey := append(PublicKey.X.Bytes(), PublicKey.Y.Bytes()...)
	return &Wallet{
		PrivateKey: PrivateKey,
		PublicKey:  publicKey,
	}
}

// 生成地址
func (self *Wallet) GenerateAddr() string {
	rip160HashValue := publicKeyHash(self.PublicKey)
	// 拼接version,总计21 bytes
	version := byte(00)
	payload := append([]byte{version}, rip160HashValue...)
	// checksum
	checkCode := CheckSum(payload)
	// 25字节数据
	payload = append(payload, checkCode...)
	// base58编码，使用btcd包:go语言实现比特币全节点源码
	address := base58.Encode(payload)
	return address
}

func publicKeyHash(data []byte) []byte {
	// 1、对公钥进行hash
	pub_hash := sha256.Sum256(data)
	// 新建一个ripmd160编码器
	hash160 := ripemd160.New()
	_, err := hash160.Write(pub_hash[:])
	if err != nil {
		log.Panic(err)
	}
	rip160HashValue := hash160.Sum(nil)
	return rip160HashValue
}

func CheckSum(payload []byte) []byte {
	// 两次256hash
	payloadHash := sha256.Sum256(payload)
	rePayloadHash := sha256.Sum256(payloadHash[:])
	checkCode := rePayloadHash[:4]
	return checkCode
}