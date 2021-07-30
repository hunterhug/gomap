# Tree Map，使用红黑树，AVL树实现

[![GitHub forks](https://img.shields.io/github/forks/hunterhug/gomap.svg?style=social&label=Forks)](https://github.com/hunterhug/gomap/network)
[![GitHub stars](https://img.shields.io/github/stars/hunterhug/gomap.svg?style=social&label=Stars)](https://github.com/hunterhug/gomap/stargazers)
[![GitHub last commit](https://img.shields.io/github/last-commit/hunterhug/gomap.svg)](https://github.com/hunterhug/gomap)
[![GitHub issues](https://img.shields.io/github/issues/hunterhug/gomap.svg)](https://github.com/hunterhug/gomap/issues)

[English README](/README.md)

哈希表在某些场景下可以称为字典，用途是可以根据 `键key` 索引该键对应的 `值value`。哈希表是什么，可以参考：[数据结构和算法（Golang实现）](https://github.com/hunterhug/goa.c)。

我们知道 `Golang map`，内置类型的 `map` 使用拉链法实现，会提前分配空间，随着元素的增加，会不断扩容，这样会一直占用空间，即使删除元素也不会缩容，导致无法垃圾回收，可能出现内存溢出的情况。

我们使用平衡二叉查找树：红黑树或者AVL树来实现 `map`，不会占用多余空间，因为 `Golang` 不支持范型，目前 `key` 必须是字符串，`value` 可以是任何类型，且为了效率，我们不进行哈希处理，也就是 `键key` 的 `hash` 就是其本身。

## 如何使用

很简单，执行：

```
go get -v github.com/hunterhug/gomap
```

目前 `gomap` 支持并发安全，并且可选多种平衡二叉查找树。

有以下几种用法:

1. `Red-Black Tree`，使用标准红黑树(2-3-4-树): `gomap.New()`，`gomap.NewMap()`，`gomap.NewRBMap()`。
2. `AVL Tree`，使用AVL树: `gomap.NewAVLMap()`。

核心 API:

```go
// Tree Map 实现
// 被设计为并发安全的
type Map interface {
	Put(key string, value interface{})            // 放入键值对
	Delete(key string)                            // 删除键
	Get(key string) (interface{}, bool)           // 获取键，返回的值 value 是 interface{} 类型的，想返回具体类型的值参考下面的方法
	GetInt(key string) (int, bool, error)         // 获取键，返回的值 value 转成 int
	GetInt64(key string) (int64, bool, error)     // 获取键，返回的值 value 转成 int64
	GetString(key string) (string, bool, error)   // 获取键，返回的值 value 转成 string
	GetFloat64(key string) (float64, bool, error) // 获取键，返回的值 value 转成 float64
	GetBytes(key string) ([]byte, bool, error)    // 获取键，返回的值 value 转成 []byte
	Contains(key string) bool                     // 查看键是否存在
	Len() int64                                   // 查看键值对数量
	KeyList() []string                            // 根据树的层次遍历，获取键列表
	KeySortedList() []string                      // 根据树的中序遍历，获取字母序排序的键列表
	Iterator() MapIterator                        // 迭代器，实现迭代
	SetComparator(comparator) Map                 // 可自定义键比较器，默认按照字母序
}

// Iterator 迭代器，不是并发安全，迭代的时候确保不会修改Map，否则可能panic或产生副作用
type MapIterator interface {
	HasNext() bool // 是否有下一对键值对
	Next() (key string, value interface{}) // 获取下一对键值对，迭代器向前一步
}
```

## 算法比较

`Red-Black Tree` 添加操作最多旋转两次，删除操作最多旋转三次，树最大高度为 `2log(N+1)`。

`AVL Tree` 添加操作最多旋转两次，删除操作旋转可能旋转到根节点，树最大高度为 `1.44log(N+2)-1.328`。

由于 `AVL Tree` 树理论上高度较低，所以查询速度稍快，但是删除操作因为旋转可能过多，会慢一点，但经过我们的优化，比如使用了非递归，只更新需要更新的树节点高度等，它们的速度是差不多的。

`Java/C#/C++` 标准库内部实现都是用 `Red-Black Tree` 实现的，`Windows` 对进程地址空间的管理则使用了 `AVL Tree`。

## 例子

下面是一个基本的例子：

```go
package main

import (
	"fmt"
	"github.com/hunterhug/gomap"
	"math/rand"
	"strconv"
	"time"
)

// loop times
// 循环的次数，用于随机添加和删除键值对
// 可以修改成1000
var num = 20

func init() {
	// random seed
	// 随机数种子初始化
	rand.Seed(time.Now().Unix())
}

// diy comparator
// 内部是按键字符串来做查找树，我们把他变成整形
func comparatorInt(key1, key2 string) int64 {
	k1, _ := strconv.Atoi(key1)
	k2, _ := strconv.Atoi(key2)
	if k1 == k2 {
		return 0
	} else if k1 > k2 {
		return 1
	} else {
		return -1
	}
}

func main() {
	checkMap := make(map[string]struct{})

	// 1. new a map default is rb tree
	// 1. 新建一个 Map，默认为标准红黑树实现
	m := gomap.New()
	//m = gomap.NewAVLMap()    // avl tree better version
	//m = gomap.NewAVLRecursionMap() // avl tree bad version

	m.SetComparator(comparatorInt) // set inner comparator

	for i := 0; i < num; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(num)))
		fmt.Println("add key:", key)
		checkMap[key] = struct{}{}

		// 2. put key pairs
		// 2. 放键值对
		m.Put(key, key)
	}

	fmt.Println("map len is ", m.Len())

	// 3. can iterator
	// 3. 迭代器使用
	iterator := m.Iterator()
	for iterator.HasNext() {
		k, v := iterator.Next()
		fmt.Printf("Iterator key:%s,value %v\n", k, v)
	}

	// 4. get key
	// 4. 获取键中的值
	key := "9"
	value, exist := m.Get(key)
	if exist {
		fmt.Printf("%s exist, value is:%s\n", key, value)
	} else {
		fmt.Printf("%s not exist\n", key)
	}

	// 5. get int will err
	// 5. 获取键中的值，并且指定类型，因为类型是字符串，所以转成整形会报错
	_, _, err := m.GetInt(key)
	if err != nil {
		fmt.Println(err.Error())
	}

	// 6. check is a rb tree
	// 6. 内部辅助，检查是不是一颗正常的红黑树
	if m.Check() {
		fmt.Println("is a rb tree,len:", m.Len())
	}

	// 7. delete '9' then find '9'
	// 7. 删除键 '9' 然后再找 '9'
	m.Delete(key)
	value, exist = m.Get(key)
	if exist {
		fmt.Printf("%s exist, value is:%s\n", key, value)
	} else {
		fmt.Printf("%s not exist\n", key)
	}

	// 8. key list
	// 8. 获取键列表
	fmt.Printf("keyList:%#v,len:%d\n", m.KeyList(), m.Len())
	fmt.Printf("keySortList:%#v,len:%d\n", m.KeySortedList(), m.Len())

	// 9. delete many
	// 9. 删除很多键值对
	for key := range checkMap {
		v, ok := m.Get(key)
		fmt.Printf("find %s:%v=%v, delete key:%s\n", key, v, ok, key)
		m.Delete(key)
		if !m.Check() {
			fmt.Println("is not a rb tree,len:", m.Len())
		}
	}

	// 10. key list
	// 10. 获取键列表
	fmt.Printf("keyList:%#v,len:%d\n", m.KeyList(), m.Len())
	fmt.Printf("keySortList:%#v,len:%d\n", m.KeySortedList(), m.Len())

	// 11. check is a rb tree
	// 11. 再次检查是否是一颗正常的红黑树
	if m.Check() {
		fmt.Println("is a rb tree,len:", m.Len())
	}
}
```

## 性能测试

我进行了压测：

```go
go test -run="bench_test.go" -test.bench=".*" -test.benchmem=1 -count=1                                                                                                            master 
goos: darwin
goarch: amd64
pkg: github.com/hunterhug/gomap
BenchmarkGolangMapPut-12                 1791264               621 ns/op             112 B/op          6 allocs/op
BenchmarkRBTMapPut-12                    1000000              1408 ns/op             104 B/op          6 allocs/op
BenchmarkAVLMapPut-12                    1000000              1440 ns/op             104 B/op          6 allocs/op
BenchmarkAVLRecursionMapPut-12           1000000              2024 ns/op             104 B/op          6 allocs/op
BenchmarkGolangMapDelete-12              3577232               303 ns/op              15 B/op          1 allocs/op
BenchmarkRBTMapDelete-12                  996924              1248 ns/op              15 B/op          1 allocs/op
BenchmarkAVLMapDelete-12                 1000000              1227 ns/op              15 B/op          1 allocs/op
BenchmarkAVLRecursionMapDelete-12         667242              1866 ns/op              15 B/op          1 allocs/op
BenchmarkGolangMapGet-12                15797131                72.2 ns/op             2 B/op          1 allocs/op
BenchmarkRBTMapGet-12                    5798295               195 ns/op               2 B/op          1 allocs/op
BenchmarkAVLMapGet-12                    5831353               197 ns/op               2 B/op          1 allocs/op
BenchmarkAVLRecursionMapGet-12           5275490               228 ns/op               2 B/op          1 allocs/op
BenchmarkGolangMapRandom-12              1256779               940 ns/op             146 B/op          8 allocs/op
BenchmarkRBTMapRandom-12                  965804              2652 ns/op             126 B/op          8 allocs/op
BenchmarkAVLMapRandom-12                  902004              2805 ns/op             126 B/op          8 allocs/op
BenchmarkAVLRecursionMapRandom-12         701880              3309 ns/op             129 B/op          8 allocs/op
PASS
ok      github.com/hunterhug/gomap      66.006s
```

其中 `GolangMap` 是内置 `map`，使用空间换时间，其速度是最快的。`RBTMap` 是非递归版本的标准红黑树，`AVLMap` 是非递归版本的 `AVL` 树，这两个查找树速度差不多，`AVLRecursionMap` 是递归版本的 `AVL` 树（效果差不要使用）。

如果对程序内存使用比较苛刻，在存储大量键值对情况下，不想浪费内存，可以使用二叉查找树实现的 `Map`。因为拉链法实现的 `golang map` 速度肯定更快，如果资源充足，直接使用官方 `map` 即可。

空间换时间，还是时间换空间，这是一个问题。

# License

```
Copyright [2019-2021] [github.com/hunterhug]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```