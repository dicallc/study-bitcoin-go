package block

import (
	"bytes"
	"wallet"
)

//输入 收账记录
type TxInput struct {
	Txid      []byte //交易ID的hash
	Vout      int    //所引用Output的索引值
	Signature []byte //数字签名
	PubKey    []byte //钱包里的公钥
}

//判断当前输入是否和某个输出吻合
func (in *TxInput) UseKey(pubKeyHash []byte) bool {
	lockingHash := wallet.HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
