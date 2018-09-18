package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

//生成公钥私钥函数
func RsaGenKey(bits int) error {
	//1.使用rsa中Generatekey方法生成私钥
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	//2. 通过X509标准将得到的RAS私钥序列化成ASN.1的DER编码字符串
	privStream := x509.MarshalPKCS1PrivateKey(privKey)
	//3. 将私钥字符串设置到pem格式块中
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privStream,
	}
	//4. 通过pem将设置好的数据进行编码，并写入磁盘文件中
	privFile, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	defer privFile.Close()
	err = pem.Encode(privFile, &block)
	if err != nil {
		return err
	}

	//1. 从得到的私钥对象中将公钥信息去除
	pubKey := privKey.PublicKey
	//2. 通过X509标准得到的rsa公钥序列化为字符串
	publicStream, err := x509.MarshalPKIXPublicKey(&pubKey)
	if err != nil {
		return err
	}
	//3. 将公钥字符串设置到pe格式块中
	block = pem.Block{
		Type:  "RSA public KEY",
		Bytes: publicStream,
	}
	//4. 通过pem将设置好的数据进行编码，并写入磁盘文件
	publicFile, err := os.Create("public.pem")
	if err != nil {
		return err
	}
	defer publicFile.Close()
	err = pem.Encode(publicFile, &block)
	if err != nil {
		return err
	}
	return nil
}

//公钥加密函数
func RsaPublicEncrypt(src []byte, pathName []byte) ([]byte, error) {
	//1.将公钥文件中的公钥读出，得到使用pem编码的字符串
	msg := []byte("")
	file, err := os.Open(string(pathName))
	if err != nil {
		return msg, err
	}
	//1.1得到文件的属性信息，通过属性信息对象得到文件大小
	info, err := file.Stat()
	if err != nil {
		return msg, err
	}
	recvbuf := make([]byte, info.Size())
	file.Read(recvbuf)
	//2.将得到的字符串解码
	block, _ := pem.Decode(recvbuf)
	//3.使用x509将编码之后的公钥解析出来
	pubInter, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return msg, err
	}
	pubKey := pubInter.(*rsa.PublicKey)
	//4.使用得到的公钥通过RSA进行数据加密
	msg, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey, src)
	if err != nil {
		return msg, err
	}
	return msg, nil
}

//使用私钥解密
func RsaPrivateDecrypt(src []byte, pathName string) ([]byte, error) {
	//1.打开私钥文件
	msg := []byte("")
	file, err := os.Open(pathName)
	if err != nil {
		return msg, err
	}
	//2.读文件内容
	info, err := file.Stat()
	recvbuf := make([]byte, info.Size())
	file.Read(recvbuf)
	//3.将得到字符串解码
	block, _ := pem.Decode(recvbuf)
	//4.通过X509还原私钥数据
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return msg, err
	}
	//5.通过私钥对数据解密
	msg, err = rsa.DecryptPKCS1v15(rand.Reader, privKey, src)
	if err != nil {
		return msg, err
	}
	return msg, nil
}
