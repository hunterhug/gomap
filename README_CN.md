# 哈希表/字典 Map 实现，底层数据结构为平衡二叉查找树

[![GitHub forks](https://img.shields.io/github/forks/hunterhug/gomap.svg?style=social&label=Forks)](https://github.com/hunterhug/gomap/network)
[![GitHub stars](https://img.shields.io/github/stars/hunterhug/gomap.svg?style=social&label=Stars)](https://github.com/hunterhug/gomap/stargazers)
[![GitHub last commit](https://img.shields.io/github/last-commit/hunterhug/gomap.svg)](https://github.com/hunterhug/gomap)
[![Go Report Card](https://goreportcard.com/badge/github.com/hunterhug/gomap)](https://goreportcard.com/report/github.com/hunterhug/gomap)
[![GitHub issues](https://img.shields.io/github/issues/hunterhug/gomap.svg)](https://github.com/hunterhug/gomap/issues)

[English README](/README.md)

哈希表在某些场景下可以称为字典，用途是可以根据 `键key` 索引该键对应的 `值value`。哈希表是什么，可以参考：[数据结构和算法（Golang实现）](https://github.com/hunterhug/goa.c)。

目前实现的哈希表 `Map`，不是用链表数组数据结构实现的，而是以平衡二叉查找树形式来实现。

我们知道 `Golang` 有标准内置类型 `map`，内置类型的 `map` 使用拉链法实现，会提前分配空间，随着元素的增加，会不断扩容，这样会一直占用空间，即使删除元素也不会缩容，导致无法垃圾回收，可能出现内存溢出的情况。

使用平衡二叉查找树结构实现的哈希表，不会占用多余空间，因为平衡树的缘故，查找元素效率也不是那么地糟糕，所以有时候选择它来做哈希表是特别好的。

## 如何使用

很简单，执行：

```
go get -v github.com/hunterhug/gomap
```

目前 `gomap` 支持并发安全，并且可选多种平衡二叉查找树。

有以下几种用法:

1. 使用标准红黑树(2-3-4-树): `gomap.New()`，`gomap.NewMap()`，`gomap.NewRBMap()`。
2. 使用AVL树: `gomap.NewAVLMap()`。

核心 API:

```go
// Map 实现
// 被设计为并发安全的
type Map interface {
	Put(key string, value interface{})            // 放入键值对
	Delete(key string)                            // 删除键
	Get(key string) (interface{}, bool)           // 获取键，返回的值value是interface{}类型的，想返回具体类型的值参考下面的方法
	GetInt(key string) (int, bool, error)         // get value auto change to Int
	GetInt64(key string) (int64, bool, error)     // get value auto change to Int64
	GetString(key string) (string, bool, error)   // get value auto change to string
	GetFloat64(key string) (float64, bool, error) // get value auto change to string
	GetBytes(key string) ([]byte, bool, error)    // get value auto change to []byte
	Contain(key string) bool                      // 查看键是否存在
	Len() int64                                   // 查看键值对数量
	KeyList() []string                            // 根据树的层次遍历，获取键列表
	KeySortedList() []string                      // 根据树的中序遍历，获取字母序排序的键列表
	Iterator() MapIterator                        // 迭代器，实现迭代
}

// Iterator 迭代器，不是并发安全，迭代的时候确保不会修改Map，否则可能panic或产生副作用
type MapIterator interface {
	HasNext() bool // 是否有下一对键值对
	Next() (key string, value interface{}) // 获取下一对键值对，迭代器向前一步
}
```

因为 `Golang` 不支持范型，目前 `key` 必须是字符串，`value` 可以是任何类型。

## 例子

下面是一个基本的例子：

```go
package main

import (
	"fmt"
	"github.com/hunterhug/gomap"
	"math/rand"
	"time"
)

// 循环的次数，用于随机添加和删除键值对
// 可以修改成1000
var num = 10

func init() {
	// 随机数种子初始化
	rand.Seed(time.Now().Unix())
}

func main() {
	// 1. 新建一个 Map，默认为标准红黑树实现
	m := gomap.New()
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(num)))
		fmt.Println("add key:", key)
		// 2. 放键值对
		m.Put(key, key)
	}

	fmt.Println("map len is ", m.Len())

	// 3. 迭代器使用
	iterator := m.Iterator()
	for iterator.HasNext() {
		k, v := iterator.Next()
		fmt.Printf("Iterator key:%s,value %v\n", k, v)
	}

	// 4. 获取键中的值
	key := "9"
	value, exist := m.Get(key)
	if exist {
		fmt.Printf("%s exist, value is:%s\n", key, value)
	} else {
		fmt.Printf("%s not exist\n", key)
	}

	// 5. 获取键中的值，并且指定类型，因为类型是字符串，所以转成整形会报错
	_, _, err := m.GetInt(key)
	if err != nil {
		fmt.Println(err.Error())
	}

	// 6. 内部辅助，检查是不是一颗正常的红黑树
	if m.Check() {
		fmt.Println("is a rb tree,len:", m.Len())
	}

	// 7. 删除键 '9' 然后再找 '9'
	m.Delete(key)
	value, exist = m.Get(key)
	if exist {
		fmt.Printf("%s exist, value is:%s\n", key, value)
	} else {
		fmt.Printf("%s not exist\n", key)
	}

	// 8. 删除很多键值对
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(num)))
		fmt.Println("delete key:", key)
		m.Delete(key)
	}

	// 9. 获取键列表
	fmt.Printf("keyList:%#v,len:%d\n", m.KeyList(), m.Len())
	fmt.Printf("keySortList:%#v,len:%d\n", m.KeySortedList(), m.Len())

	// 10. 再次检查是否是一颗正常的红黑树
	if m.Check() {
		fmt.Println("is a rb tree,len:", m.Len())
	}
}
```

## 性能测试

```go
go test -run="bench_test.go" -test.bench=".*" -test.benchmem=1 -count=1

BenchmarkBRTMapPut-4             1000000              2076 ns/op              95 B/op          3 allocs/op
BenchmarkGolangMapPut-4          1580353               696 ns/op             139 B/op          3 allocs/op
BenchmarkBRTMapDelete-4          1000000              1134 ns/op              16 B/op          2 allocs/op
BenchmarkGolangMapDelete-4       4215994               293 ns/op              16 B/op          2 allocs/op
BenchmarkBRTMapRandom-4           675942              3547 ns/op             163 B/op          8 allocs/op
BenchmarkGolangMapRandom-4        650377              1720 ns/op             225 B/op          8 allocs/op
```

如果对程序内存空间的占用要求比较高，在存储大量键值对情况下，不想浪费内存，可以使用二叉查找树实现的 `Map`。

因为拉链法实现的 `golang map` 速度肯定更快，如果资源充足，直接使用官方 `map` 即可。

空间换时间，还是时间换空间，这是一个问题。
