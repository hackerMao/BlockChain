package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type Person struct {
	Name string
	Age uint
}

func main() {
	// 编码
	p := Person{
		Name:"xiaoMing",
		Age: 20,
	}
	// 新建一个编码器
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	// 使用编码器进行编码
	err := encoder.Encode(&p)
	if err != nil {
		log.Panic("编码出错了")
	}
	fmt.Printf("编码结果：%v\n", buffer.Bytes())
	// 解码
	var P Person
	// 新建一个解码器
	decoder := gob.NewDecoder(bytes.NewReader(buffer.Bytes()))
	// 使用解码器解码
	decoder.Decode(&P)
	fmt.Printf("解码结果：%v\n", P)
}
