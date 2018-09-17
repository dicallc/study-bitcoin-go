package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"utils"
)

//填充最后一个分组的函数
//src-原始数据
//blockSize-每一个分组的数据长度
func paddingText(src []byte, blockSize int) []byte {
	//1.求出最后一个分组要填充多少个字节
	padding := blockSize - len(src)%blockSize
	//2.创建新的切片，切片的字节数为padding，并初始化
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	//3.将创建出的新切片和原始数据进行连接
	return append(src, padText...)
}

//删除末尾填充的字节
func unPaddingText(src []byte) []byte {
	//1.求出要处理的切片的长度
	lens := len(src)
	//2.取出最后一个字符，得到其整型值
	number := int(src[lens-1])
	//3.将切片末尾的number个字节删除
	newText := src[:lens-number]
	return newText
}

//使用des进行对称加密
func encryptDES(src, key []byte) []byte {
	//1.创建并返回一个使用DES算法cipher.block接口
	block, err := des.NewCipher(key)
	utils.CheckErr("encryptDES", err)
	//2.对最后一个明文分组进行数据填充
	src = paddingText(src, block.BlockSize())
	//3.创建一个密码分组为链接模式，底层使用DES加密的BlockMode接口
	iv := []byte("aaaabbbb")
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//4.加密连续的数据块
	dst := make([]byte, len(src))
	blockMode.CryptBlocks(dst, src)
	return dst
}

//使用3des加密
func encrypt3DES(src, key []byte) []byte {
	//1.创建并返回一个使用3DES算法cipher.block接口
	block, err := des.NewTripleDESCipher(key)
	utils.CheckErr("encrypt3DES", err)
	//2.对最后一个明文分组进行数据填充
	src = paddingText(src, block.BlockSize())
	//3.创建一个密码分组为链接模式，底层使用3DES加密的BlockMode接口
	blockMode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	//4.加密连续的数据块 数据覆盖
	blockMode.CryptBlocks(src, src)
	return src
}
func encryptAES(src, key []byte) []byte {
	//1.创建并返回一个使用3DES算法cipher.block接口
	block, err := aes.NewCipher(key)
	utils.CheckErr("encryptAES", err)
	//2.对最后一个明文分组进行数据填充
	src = paddingText(src, block.BlockSize())
	//3.创建一个密码分组为链接模式，底层使用3DES加密的BlockMode接口
	blockMode := cipher.NewCBCEncrypter(block, key)
	//4.加密连续的数据块 数据覆盖
	blockMode.CryptBlocks(src, src)
	return src
}

//使用des解密
func decryptDES(src, key []byte) []byte {
	//1.创建并返回一个使用DES算法的cipher.Block接口
	block, err := des.NewCipher(key)
	utils.CheckErr("decryptDES", err)
	//2.创建一个密码分组为链接模式的，底层使用DES解密的BlockMode接口
	iv := []byte("aaaabbbb")
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//3.数据块解密
	blockMode.CryptBlocks(src, src)
	//4.去掉最后一组的填充数据
	return unPaddingText(src)

}

//使用3des对数据解密

func decrypt3DES(src, key []byte) []byte {
	//1.创建并返回一个使用DES算法的cipher.Block接口
	block, err := des.NewTripleDESCipher(key)
	utils.CheckErr("decrypt3DES", err)
	//2.创建一个密码分组为链接模式的，底层使用DES解密的BlockMode接口
	iv := []byte("aaaabbbb")
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//3.数据块解密
	blockMode.CryptBlocks(src, src)
	//4.去掉最后一组的填充数据
	return unPaddingText(src)
}
func decryptAES(src, key []byte) []byte {
	//1.创建并返回一个使用DES算法的cipher.Block接口
	block, err := aes.NewCipher(key)
	utils.CheckErr("decryptAES", err)
	//2.创建一个密码分组为链接模式的，底层使用DES解密的BlockMode接口
	blockMode := cipher.NewCBCDecrypter(block, key)
	//3.数据块解密
	blockMode.CryptBlocks(src, src)
	//4.去掉最后一组的填充数据
	return unPaddingText(src)
}
