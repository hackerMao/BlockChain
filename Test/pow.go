package Test

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	data := "helloworld"
	for i:=0; i<1000000; i++ {
		hashText := sha256.Sum256([]byte(data+string(i)))
		fmt.Printf("hash: %x\n", string(hashText[:]))
	}
}
