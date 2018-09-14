master:引入了pow工作证明
源码有分支:
db:引入了数据库进行存储

## db版本的思路：

### 1.blot数据的介绍

 key->value进行读取存储，轻量级，开源的

### 2.NewBlockChain函数的重写

有数组编程操作数据库

创建数据库文件

### 3.addBlock的函数重写

对数据的读取和写入

### 4.对数据库的遍历

迭代器的编写，Iterator

### 5.命令行介绍及编写

a.添加区块命令

b.打印区块链命令

## utxo版本的思路：

添加交易utxo，utxo创建，转移等复杂工作

### 1.创建区块链的创世区块操作放到命令

### 2.定义交易结构
    - 交易id
    - 交易输出
    - 交易输入

### 3.根据交易结构改写代码
    a.创建区块链的时候生成奖励
    b.通过指定地址检索到他们相关的utxo
    c.实现utxo的转移(创建交易函数:NewTransaction(from, tostring ,amount float64))
go env