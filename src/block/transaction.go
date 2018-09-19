package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"utils"
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
	Signature []byte //签名
	PubKey    []byte //公钥
	//ScriptSig string //解锁脚本
}

//一个事物输出
type TXOutput struct {
	Value int //支付给收款方金额值
	//ScriptPubKey string //锁定脚本，指定收款方的地址
	PubKeyHash []byte //解锁脚本key
}

//检查输入是否使用特定的键来解锁输出
func (in *TXInput) UseKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

//lock只需锁定输出
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := utils.Base58Encode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

//检查提供散列是否用于锁定输出
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

//新的交易输出
func NewTxOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
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
	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
	//subsidy是奖励的金额
	txOut := NewTxOutput(subsidy, to)
	//txout := TXOutput{subsidy, }
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txOut}}
	tx.SetID()
	return &tx
}

//签名
func (tx *Transaction) sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	//判断是否是coinbase交易
	if tx.IsCoinbase() {
		return
	}
	for _, vin := range tx.TxInputs {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	txCopy := tx.TrimmedCopy()
	for inID, vin := range txCopy.TxInputs {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.TxInputs[inID].Signature = nil
		txCopy.TxInputs[inID].PubKey = prevTx.TXOutputs[vin.Vout].PubKeyHash
		txCopy.SetID()
		txCopy.TxInputs[inID].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.TxInputs[inID].Signature = signature
	}

}

//创建一个副本用于签名
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.TxInputs {
		//该副本将包括所有的输入和输出，但TXInput.Signature并TXInput.PubKey设置为零
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout, nil, nil})
	}
	for _, vout := range tx.TXOutputs {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

//检查当前的用户能否解开引用的utxo
//func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
//	return in.ScriptSig == unlockingData
//}

//检查当前用户时候的utxo的所有者
//func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
//	return out.ScriptPubKey == unlockingData
//}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.TxInputs) == 1 && len(tx.TxInputs[0].Txid) == 0 && tx.TxInputs[0].Vout == -1
}

//创建普通交易，send的辅助函数
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	//返回合适的UTXO
	acc, validOutputs := bc.FindSuitableUTXOs(from, amount)
	//判断是否有那么多可花费的币
	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}
	//遍历有效UTXO的合集
	for txid, outIndexs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		CheckErr(err)
		//遍历所有引用UTXO的索引，每一个索引需要创建一个Input
		for _, outindex := range outIndexs {
			input := TXInput{txID, outindex, from}
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
