[ripemd160](https://blog.csdn.net/u013397318/article/details/80937583)

[golang利用gob序列化struct对象保存到本地](https://studygolang.com/articles/2888)

### 公钥私钥转换过程

![mark](http://7xnk07.com1.z0.glb.clouddn.com/blog/180920/H2HIHaE78d.png?imageslim)

### 序列化内容代码

```
	var content bytes.Buffer
	encoder:=gob.NewEncoder(&content)
	err:=encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
```

#### 增强序列化内容gob

```
gob.Register(elliptic.P256())//注册的目的是为了可以序列化任何类型
```



### 文件操作:ioutil

```
err=ioutil.WriteFile(walletFile,content.Bytes(),0644) //0644是文件权限
```



### 使用脚本解锁区块

之前只是对比一下用户名是否一样，现在必须使用其他值去代替

- 数字签名
- 公钥

### pri-public

就在于小写和大写第一个字母区别

### 重构过程

1. 引进钱包概念
2. 修改UTXO的输入输出-主要是重构之前的用户名验证变成数字签名和公钥
3. 正式去修改UTXO

