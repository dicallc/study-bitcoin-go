package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

//使用md5对数据进行哈希运算
func GetMd5str_1(src []byte) string {
	//1.给哈希算法添加数据
	res := md5.Sum(src)
	//2.将数据格式化成16进制
	//myres:=fmt.Sprintf("%x",res)
	myres := hex.EncodeToString(res[:])
	return myres
}
func GetMd5str_2(src []byte) string {
	//1.给哈希算法添加数据
	myHash := md5.New()
	//2.将数据格式化成16进制
	io.WriteString(myHash, string(src))
	//myHash.Write(src)
	res := myHash.Sum(nil)
	//4.散列值格式化
	return hex.EncodeToString(res)
}
