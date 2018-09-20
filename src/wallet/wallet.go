package wallet

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"utils"
)

const version = byte(0x00) //16进制 0 版本号
const walletFile = "db/wallet.dat"
const addressChecksumlen = 4 //地址检查长度4

//钱包

type Wallet struct {
	/**
	PrivateKey: ECDSA基于椭圆曲线
	使用曲线生成私钥，并从私钥生成公钥
	*/
	PrivateKey ecdsa.PrivateKey //私钥
	PublicKey  []byte           //公钥
}

//创建一个新钱包
func NewWallet() *Wallet {
	//公钥私钥生成
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

//得到一个钱包地址

func (w Wallet) GetAddress() []byte {
	//1.使用 RIPEMD160(SHA256(PubKey)) 哈希算法，取公钥并对其哈希两次
	pubKeyHash := HashPubKey(w.PublicKey)
	//2.给哈希加上地址生成算法版本的前缀
	versionedPayload := append([]byte{version}, pubKeyHash...)
	//3.对于第二步生成的结果，使用 SHA256(SHA256(payload)) 再哈希，计算校验和。校验和是结果哈希的前四个字节
	checksum := checksum(versionedPayload)
	//4.将校验和附加到 version+PubKeyHash 的组合中
	fullPayload := append(versionedPayload, checksum...)
	//5.使用 Base58 对 version+PubKeyHash+checksum 组合进行编码
	address := utils.Base58Encode(fullPayload)
	return address

}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:addressChecksumlen]
}

// 使用RIPEMD160(SHA256(PubKey))哈希算法得到Hashpubkey
func HashPubKey(pubKey []byte) []byte {
	//1.256
	publicSHA256 := sha256.Sum256(pubKey)
	//2.160
	RIPEMD160Hasher := crypto.RIPEMD160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	utils.CheckErr("", err)
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

//椭圆算法返回私钥与公钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	//实现了P-256的曲线
	curve := elliptic.P256()
	//获取私钥
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	utils.CheckErr("", err)
	//在基于椭圆曲线的算法中，公钥是曲线上的点。因此，公钥是X，Y坐标的组合。在比特币中，这些坐标被连接起来形成一个公钥
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}
