package block

import (
	"bytes"
	"encoding/gob"
	"time"
)

//区块结构
type Block struct {
	Hash          []byte         //hase值
	Transactions  []*Transaction //交易数据
	PrevBlockHash []byte         //存储前一个区块的Hase值
	Timestamp     int64          //生成区块的时间
	Nonce         int            //工作量证明算法的计数器
	Height        int
}

//生成一个新的区块方法
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int) *Block {
	//GO语言给Block赋值{}里面属性顺序可以打乱，但必须制定元素 如{Timestamp:time.Now().Unix()...}
	block := &Block{Timestamp: time.Now().Unix(), Transactions: transactions, PrevBlockHash: prevBlockHash, Hash: []byte{}, Nonce: 0, Height: height}

	//工作证明
	pow := NewProofOfWork(block)
	//工作量证明返回计数器和hash
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

//序列化Block
func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	CheckErr(err)
	return result.Bytes()
}

//反序列化
func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	CheckErr(err)
	return &block
}

//区块校验
func (i *Block) Validate() bool {
	return NewProofOfWork(i).Validate()
}

//需要将Txs转换成[]byte
func (b *Block) HashTransactions() []byte {
	var transactions [][]byte
	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := NewMerkleTree(transactions)
	return mTree.RootNode.Data
}

//创世块方法
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{}, 0)
}
