package main

import (
	"fmt"
	"utils"
)

func main() {
	bytes := []byte("http://liyuechun.org")

	bytes58 := utils.Base58Encode(bytes)

	fmt.Printf("%x\n", bytes58)

	fmt.Printf("%s\n", bytes58)

	bytesStr := utils.Base58Decode(bytes58)

	fmt.Printf("%s\n", bytesStr[1:])
}
