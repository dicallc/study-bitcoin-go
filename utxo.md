# UTXO

前言：

这一节内容我自己也弄混了，很晕，看代码，找资料，看视频，都力不从心

可能是确实比较难

点出一点让有幸者看到的人不晕

1. 区块链里面包含多个区块
2. 区块包含很多信息其中就有UTXO(Transaction[])
3. 而很多UTXO就是一个账本包含了该区块的交易信息

觉得有点晕就默读几遍

### 1.UTXO是什么

我看了很多文章，写的很专业，专业的让人看不懂

[]: http://8btc.com/article-4381-1.html	"其实并没有什么比特币，只有 UTXO"



其实吧，把UTXO比喻为一本账本是很容易理解的，但是你得明白UTXO绝不是账本，他比账本牛逼，从上面文章可以知道答案

### 2.UTXO的结构

```
type Transaction struct {
	ID        []byte
	TxInputs  []TXInput  //输入
	TXOutputs []TXOutput //输出
}
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
```

一本账本有进有出则如代码所示

#### TXInput:

指明交易发起人可支付资金的来源，包含：

* 引用utxo所在交易id
* 所消费utxo在output中的索引
* 解锁脚本

#### TXOutput

包含资金接收方的相关信息，包含：

* 接收金额
* 锁定脚本



接下来从上面知道的概念，去理解UTXO消耗过程

* lily给Alice转账
* Jim给Alice转账
* Alice给bob转账



![mark](http://7xnk07.com1.z0.glb.clouddn.com/blog/180915/IjjIEl9B80.png?imageslim)



这张图，主要用来理解UTXO的结构，以及生成消耗过程

#### coinbase交易

coinbase交易是一种特殊类型的交易，不需要以前存在的输出。它无处不在地创造产出（即“硬币”）。没有鸡的鸡蛋。这是矿工获得开采新矿区的奖励。 

我们就先从简单的创世区块开始理解UTXO：

```
const subsidy = 10 //初始化奖励为10
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
```

可以看到TXInput，因为没有交易ID的hash，所引用Output的索引值，所以input只需要传解锁脚本即可，而TxOut的奖励会随着以后的难度进行调整

SetID：方法和之前序列化基本一样

#### 修改创建区块也要引入UTXO

之前我在前言就说了，区块中包含UTXO，所以结构就要改了，每个区块都有一个账本就是UTXO

```
type Block struct {
	Hash          []byte         //hase值
	Transactions  []*Transaction //交易数据
	PrevBlockHash []byte         //存储前一个区块的Hase值
	Timestamp     int64          //生成区块的时间
	Nonce         int            //工作量证明算法的计数器
}
//创世块方法 这个Transaction 就是之前写的NewCoinbaseTX可获得
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}
// 创新一个新的区块数据
func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}
		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bc := Blockchain{tip, db}
	return &bc
}

```

同时工作量证明中也要修改，引入UTXO进行验证



```
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	//注意一定要将原始数据转换成[]byte，不能直接从字符串转
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			utils.IntToHex(pow.block.Timestamp),
			utils.IntToHex(int64(targetBits)),
			utils.IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}
```

而验证UTXO的办法也就是把账本中所有UTXO的id，连接起来进行哈希值计算



### 3.测试创建区块链指令



cli代码中引入指令，主要看实际功能性代码

```
createblockchain -address dicallc//这个是区块名字之后引用为地址

case createblockchain:
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		
......
if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}
	
func (cli *CLI) createBlockchain(address string) {
	bc := block.CreateBlockchain(address)//该方法就是之前提到的
	block.Close(bc) //关闭数据库
	fmt.Println("Done!")
}
```

### 4.查询余额指令

既然你得到的经历，那我们应该可以去查询余额

```
//查询余额
func (cli *CLI) getBalance(address string) {
	bc := block.GetBlockChainHandler()//获取区块链对象和数据库对象
	defer block.Close(bc)
	balance := 0
	//查询所有未经使用的交易地址
	UTXOS := bc.FindUTXO(address)
	//算出未使用的交易地址的value
	for _, out := range UTXOS {
		balance += out.Value
	}
	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
```

这里涉及到一个很重要的方法

#### FindUTXO

思路：找到所有能够使用UTXO。再从其中找到指定地址

其中有个逻辑：

> 如果一区块Utxo的Input解锁地址和一个区块Out锁定地址一样了就是已经消费掉了

```
//寻找指定地址能够使用的utxo
func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var Utxos []TXOutput
	//未使用的UTXO
	unspentTransactions := bc.FindUnspentTransactions(address)
	//遍历交易
	for _, tx := range unspentTransactions {
		//遍历output
		for _, out := range tx.TXOutputs {
			//当前地址拥有的utxo
			if out.CanBeUnlockedWith(address) {
				Utxos = append(Utxos, out)
			}
		}
	}
	return Utxos
}

//返回指定地址能够支配的utxo的交易集合
func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction //
	spentTXOs := make(map[string][]int)
	//迭代所有区块
	bci := bc.Iterator()
	for {
		block := bci.Next()
		//从区块中取出账本迭代UTXO
		for _, tx := range block.Transactions {
			txId := hex.EncodeToString(tx.ID) //交易id转成String
			//去除挖矿奖励
			if tx.IsCoinbase() == false {
				//目的：找到已经消耗的UTXO，把他们放到一个集合里
				//需要两个字段来表示使用过的utxo:a.交易id,b.output的缩影
				for _, input := range tx.TxInputs {
					if input.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(input.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], input.Vout)
					}
				}
			}
			//遍历output 目的找到所有能支配的utxo 这里有个逻辑如果一区块Utxo的Input解锁地址和一个区块Out锁定地址一样了就是已经消费掉了
		Outputs:
			for outIdx, out := range tx.TXOutputs {
				if spentTXOs[txId] != nil {
					for _, spentOut := range spentTXOs[txId] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

		}
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
}
```

### 5.转币的指令

现在，我们要发送一些硬币给别人。为此，我们需要创建一个新的UTXO，将它放在一个区块中，然后挖掘块

思路：这里肯定是需要根据地址去寻找区块集合，然后去凑齐金额

```
//发送
func (cli *CLI) send(from, to string, amount int) {
	bc := block.GetBlockChainHandler()
	defer block.Close(bc)

	tx := block.NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*block.Transaction{tx})
	fmt.Println("Success!")
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

func (bc *Blockchain) FindSuitableUTXOs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0
Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.TXOutputs {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOutputs
}
```

 