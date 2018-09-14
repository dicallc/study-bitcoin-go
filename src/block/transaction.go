package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

//交易事务
type Transaction struct {
	ID        []byte
	TxInputs  []TXInput  //输入
	TXOutputs []TXOutput //输出
}

//指明交易发起人可支付资金的来源
type TXInput struct {
	Txid      []byte //交易ID的hash
	Vout      int    //所引用Output的索引值
	ScriptSig string //解锁脚本
}

//一个事物输出
type TXOutput struct {
	Value        int    //支付给收款方金额值
	ScriptPubKey string //锁定脚本，指定收款方的地址
}

//设置交易ID hash
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte //32位的hash字节
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	CheckErr(err)
	//将交易信息sha256
	hash = sha256.Sum256(encoded.Bytes())
	//生成hash
	tx.ID = hash[:]
}

const subsidy = 10 //初始化补助为10
//创建Coinbase 只有收款人，没有付款人，是矿工的奖励交易
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'")
	}
	txin := TXInput{[]byte{}, -1, data}
	//subsidy是奖励的金额
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()
	return &tx
}

//检查当前的用户能否解开引用的utxo
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

//检查当前用户时候的utxo的所有者
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.TxInputs) == 1 && len(tx.TxInputs[0].Txid) == 0 && tx.TxInputs[0].Vout == -1
}

//创建普通交易，send的辅助函数
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	//返回合适的UTXO
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	//判断是否有那么多可花费的币
	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}
	//1.创建inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}
	// Build a list of outputs
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from}) // a change
	}
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx

}
