# Map Golang implement by Red-Black Tree, AVL Tree

[![GitHub forks](https://img.shields.io/github/forks/hunterhug/gomap.svg?style=social&label=Forks)](https://github.com/hunterhug/gomap/network)
[![GitHub stars](https://img.shields.io/github/stars/hunterhug/gomap.svg?style=social&label=Stars)](https://github.com/hunterhug/gomap/stargazers)
[![GitHub last commit](https://img.shields.io/github/last-commit/hunterhug/gomap.svg)](https://github.com/hunterhug/gomap)
[![Go Report Card](https://goreportcard.com/badge/github.com/hunterhug/gomap)](https://goreportcard.com/report/github.com/hunterhug/gomap)
[![GitHub issues](https://img.shields.io/github/issues/hunterhug/gomap.svg)](https://github.com/hunterhug/gomap/issues)

[中文说明](/README.md)

Map implement by tree data struct such Red-Black Tree, AVL Tree.

Our tree map is design to be concurrent safe, and don't waste any space different from golang standard map which not shrink even map key pairs num is 0.

## Usage

Simple get it by:

```
go get -v github.com/hunterhug/gomap
```

There are:

1. Standard Red-Black Tree Map(2-3-4-Tree): `gomap.New()`，`gomap.NewMap()`,`gomap.NewRBMap()`.
2. AVL Tree Map: `gomap.NewAVLMap()`.

Core api:

```go
// Map method
// design to be concurrent safe
type Map interface {
	Put(key string, value interface{})                            // put key pairs
	Delete(key string)                                            // delete a key
	Get(key string) (value interface{}, exist bool)               // get value from key
	GetInt(key string) (value int, exist bool, err error)         // get value auto change to Int
	GetInt64(key string) (value int64, exist bool, err error)     // get value auto change to Int64
	GetString(key string) (value string, exist bool, err error)   // get value auto change to string
	GetFloat64(key string) (value float64, exist bool, err error) // get value auto change to string
	GetBytes(key string) (value []byte, exist bool, err error)    // get value auto change to []byte
	Contains(key string) (exist bool)                             // map contains key?
	Len() int64                                                   // map key pairs num
	KeyList() []string                                            // map key out to list from top to bottom which is layer order
	KeySortedList() []string                                      // map key out to list sorted
	Iterator() MapIterator                                        // map iterator, iterator from top to bottom which is layer order
	MaxKey() (key string, value interface{}, exist bool)          // find max key pairs
	MinKey() (key string, value interface{}, exist bool)          // find min key pairs
	SetComparator(comparator) Map                                 // set compare func to control key compare
}

// Iterator concurrent not safe
// you should deal by yourself
type MapIterator interface {
	HasNext() bool
	Next() (key string, value interface{})
}
```

We has already implement them by non recursion way and optimized a lot, so use which type of tree map is no different.

## Example

Some example below:

```go
package main

import (
	"fmt"
	"github.com/hunterhug/gomap"
	"math/rand"
	"time"
)

// loop times
var num = 10

func init() {
	// random seed
	rand.Seed(time.Now().Unix())
}

func main() {
	// 1. new a map
	m := gomap.New()
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(num)))
		fmt.Println("add key:", key)
		// 2. put key pairs
		m.Put(key, key)
	}

	fmt.Println("map len is ", m.Len())

	// 3. can iterator
	iterator := m.Iterator()
	for iterator.HasNext() {
		k, v := iterator.Next()
		fmt.Printf("Iterator key:%s,value %v\n", k, v)
	}

	// 4. get key
	key := "9"
	value, exist := m.Get(key)
	if exist {
		fmt.Printf("%s exist, value is:%s\n", key, value)
	} else {
		fmt.Printf("%s not exist\n", key)
	}

	// 5. get int will err
	_, _, err := m.GetInt(key)
	if err != nil {
		fmt.Println(err.Error())
	}

	// 6. check is a rb tree
	if m.Check() {
		fmt.Println("is a rb tree,len:", m.Len())
	}

	// 7. delete '9' then find '9'
	m.Delete(key)
	value, exist = m.Get(key)
	if exist {
		fmt.Printf("%s exist, value is:%s\n", key, value)
	} else {
		fmt.Printf("%s not exist\n", key)
	}

	// 8. delete many
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(num)))
		fmt.Println("delete key:", key)
		m.Delete(key)
	}

	// 9. key list
	fmt.Printf("keyList:%#v,len:%d\n", m.KeyList(), m.Len())
	fmt.Printf("keySortList:%#v,len:%d\n", m.KeySortedList(), m.Len())

	// 10. check is a rb tree
	if m.Check() {
		fmt.Println("is a rb tree,len:", m.Len())
	}
}
```

## BenchTest

We test `Golang map` and `Red-Black Tree`, `AVL Tree`:

```go
go test -run="bench_test.go" -test.bench=".*" -test.benchmem=1 -count=1

BenchmarkGolangMapPut-4                  1000000              1385 ns/op             145 B/op          6 allocs/op
BenchmarkRBTMapPut-4                      528231              3498 ns/op             113 B/op          6 allocs/op
BenchmarkAVLMapPut-4                     1000000              3317 ns/op             104 B/op          6 allocs/op
BenchmarkAVLRecursionMapPut-4             389806              4563 ns/op             116 B/op          6 allocs/op
BenchmarkGolangMapDelete-4               2630281               582 ns/op              15 B/op          1 allocs/op
BenchmarkRBTMapDelete-4                  2127256               624 ns/op              15 B/op          1 allocs/op
BenchmarkAVLMapDelete-4                   638918              2256 ns/op              15 B/op          1 allocs/op
BenchmarkAVLRecursionMapDelete-4          376202              2813 ns/op              15 B/op          1 allocs/op
BenchmarkGolangMapGet-4                  9768266               172 ns/op               2 B/op          1 allocs/op
BenchmarkRBTMapGet-4                     3276406               352 ns/op               2 B/op          1 allocs/op
BenchmarkAVLMapGet-4                     3724939               315 ns/op               2 B/op          1 allocs/op
BenchmarkAVLRecursionMapGet-4            2550055               462 ns/op               2 B/op          1 allocs/op
BenchmarkGolangMapRandom-4               1000000              2292 ns/op             163 B/op          8 allocs/op
BenchmarkRBTMapRandom-4                   244311              4635 ns/op             136 B/op          8 allocs/op
BenchmarkAVLMapRandom-4                   488001              5879 ns/op             132 B/op          8 allocs/op
BenchmarkAVLRecursionMapRandom-4          211246              5411 ns/op             138 B/op          8 allocs/op
```

If you want to save memory space, you can choose our tree map.