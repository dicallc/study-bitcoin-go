package block

import (
	"bytes"
	"utils"
)

//输出 付钱记录
type TXOutput struct {
	Value      int    //支付给收款方金额值
	PubKeyHash []byte //解锁脚本key
}

//lock只需锁定输出 人话版本：你给别人钱你的把别人地址给设置进入生成一个输出
func (out *TXOutput) Lock(address []byte) {
	//人话：最后的地址 (Base58Decode)解密 ,然后减去后四位的版本信息就得到 pubKeyHash
	pubKeyHash := utils.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

//检查提供的公钥散列是否用于锁定输出
func (out *TXOutput) IsLockWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

//新的交易输出 人话转账给别人是要锁定一下
func NewTxOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}
