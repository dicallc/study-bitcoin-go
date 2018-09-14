# UTXO

前言：

这一节内容我自己也弄混了，很晕，看代码，找资料，看视频，都力不从心

可能是确实比较难

点出一点让有幸者看到的人不晕

1. 区块链里面包含多个区块
2. 区块包含很多信息其中就有UTXO(Transaction[])
3. 而UTXO就是一个账本包含了该区块的交易信息

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