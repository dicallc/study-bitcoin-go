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
	"wallet"
)

//交易事务
type Transaction struct {
	ID        []byte
	TxInputs  []TxInput  //输入
	TXOutputs []TXOutput //输出

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
	txin := TxInput{[]byte{}, -1, nil, []byte(data)}
	//subsidy是奖励的金额
	txOut := NewTxOutput(subsidy, to)
	//txout := TXOutput{subsidy, }
	tx := Transaction{nil, []TxInput{txin}, []TXOutput{*txOut}}
	tx.SetID()
	return &tx
}

//签名过程
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	//判断是否是coinbase交易
	if tx.IsCoinbase() {
		return
	}
	//TX 也是4
	for _, vin := range tx.TxInputs {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	//Block4 UTXO 4444也拷贝一份
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

//创建一个副本用于签名 只是Input发生了变化 公钥没了
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TXOutput

	for _, vin := range tx.TxInputs {
		//该副本将包括所有的输入和输出，但TXInput.Signature并TXInput.PubKey设置为零
		inputs = append(inputs, TxInput{vin.Txid, vin.Vout, nil, nil})
	}
	for _, vout := range tx.TXOutputs {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.TxInputs) == 1 && len(tx.TxInputs[0].Txid) == 0 && tx.TxInputs[0].Vout == -1
}

//创建普通交易，send的辅助函数
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TXOutput
	wallets, err := wallet.NewWallets()
	CheckErr(err)
	part_wallet := wallets.GetWallet(from)
	pubKeyHash := wallet.HashPubKey(part_wallet.PublicKey)

	//返回合适的UTXO
	acc, validOutputs := bc.FindSuitableUTXOs(pubKeyHash, amount)
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
			input := TxInput{txID, outindex, nil, pubKeyHash}
			inputs = append(inputs, input)
		}
	}
	// Build a list of outputs
	outputs = append(outputs, *NewTxOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTxOutput(acc-amount, from)) // a change
	}
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	//签名
	bc.SignTransaction(&tx, part_wallet.PrivateKey)
	return &tx

}
