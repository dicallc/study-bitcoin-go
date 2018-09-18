package main

import "fmt"

func main() {
	//cli.Start()
	//aesTest()
	//RsaTest()
	HashTest()
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

//测试RSA的加解密测试
func RsaTest() {
	err := RsaGenKey(4096)
	fmt.Println("错误信息：", err)
	//加密
	src := []byte("少壮不努力，老大徒伤悲")
	data, err := RsaPublicEncrypt(src, []byte("public.pem"))
	//解密
	data, err = RsaPrivateDecrypt(data, "private.pem")
	fmt.Println("解密之后的明文： " + string(data))
}

//哈希算法测试
func HashTest() {
	src := []byte("少壮不努力，老大徒伤悲")
	fmt.Println("解密之后的明文： " + GetMd5str_1(src))
	fmt.Println("解密之后的明文： " + GetMd5str_2(src))
}
