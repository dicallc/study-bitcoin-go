package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
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
	//TX 也是4  根据输入去找有不有对应的输出
	for _, vin := range tx.TxInputs {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	//txCopy Block4 UTXO      4444也拷贝一份 prevTXs:Block2 2222
	txCopy := tx.TrimmedCopy()
	//这里for循环就是把制作UTXO的输入 重新做了一遍 防止篡改？
	for inID, vin := range txCopy.TxInputs {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.TxInputs[inID].Signature = nil
		//所消费的Out公钥被引进作为了pubKey
		txCopy.TxInputs[inID].PubKey = prevTx.TXOutputs[vin.Vout].PubKeyHash
		//根据UTXO 数据生成id
		txCopy.ID = txCopy.Hash()
		txCopy.TxInputs[inID].PubKey = nil
		//真正的签名部分 钱包-私钥，UTXO-ID
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.TxInputs[inID].Signature = signature
	}
}

// ID 就是输入和输出序列化后 sha256的值
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}
	//输入 输出 id
	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}
func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
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
func DeserializeTransaction(data []byte) Transaction {
	var transaction Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}

	return transaction
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.TxInputs) == 1 && len(tx.TxInputs[0].Txid) == 0 && tx.TxInputs[0].Vout == -1
}

//创建普通交易，send的辅助函数
func NewUTXOTransaction(mWallet *wallet.Wallet, to string, amount int, UTXOSet *UTXOSet) *Transaction {
	var inputs []TxInput
	var outputs []TXOutput

	pubKeyHash := wallet.HashPubKey(mWallet.PublicKey)
	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TxInput{txID, out, nil, mWallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	from := fmt.Sprintf("%s", mWallet.GetAddress())
	outputs = append(outputs, *NewTxOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTxOutput(acc-amount, from)) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	UTXOSet.Blockchain.SignTransaction(&tx, mWallet.PrivateKey)

	return &tx
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		log.Panic("IsCoinbase")
		return true
	}
	for _, vin := range tx.TxInputs {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()
	for inID, vin := range tx.TxInputs {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.TxInputs[inID].Signature = nil
		txCopy.TxInputs[inID].PubKey = prevTx.TXOutputs[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.TxInputs[inID].PubKey = nil
		//私钥 ID
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])
		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		//公钥和id，以及签名数据求证
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}
	return true
}
