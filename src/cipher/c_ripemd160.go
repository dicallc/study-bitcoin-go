package main

import (
	"fmt"
	"golang.org/x/crypto/ripemd160"
)

func main() {
	hasher := ripemd160.New()
	hasher.Write([]byte("www.baidu.com"))
	//转成了byte
	hashBytes := hasher.Sum(nil)
	//转成hash
	//mdStr := hex.EncodeToString(md)
	hashString := fmt.Sprintf("%x", hashBytes)
	fmt.Println(hashString)
}
