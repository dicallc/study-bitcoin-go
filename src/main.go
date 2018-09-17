package main

import "fmt"

func main() {
	//cli.Start()
	aesTest()
}

//测试des的加解密
func desTest() {
	fmt.Println("========des 加解密======")
	src := []byte("少壮不努力，老大徒伤悲")
	key := []byte("12345678")
	enStr := encryptDES(src, key)
	deStr := decryptDES(enStr, key)
	fmt.Println("解密之后的明文： " + string(deStr))
}
func aesTest() {
	fmt.Println("========aes 加解密======")
	src := []byte("少壮不努力，老大徒伤悲")
	key := []byte("12345678asdfghjk")
	enStr := encryptAES(src, key)
	deStr := decryptAES(enStr, key)
	fmt.Println("解密之后的明文： " + string(deStr))
}
