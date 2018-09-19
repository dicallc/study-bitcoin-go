package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	hasher := sha256.New()
	hasher.Write([]byte("www.baidu.com"))
	//转成了byte
	hashBytes := hasher.Sum(nil)
	//转成hash
	//mdStr := hex.EncodeToString(md)
	hashString := fmt.Sprintf("%x", hashBytes)
	fmt.Println(hashString)
}
