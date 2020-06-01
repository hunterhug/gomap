/*
	All right reserved：https://github.com/hunterhug/gomap at 2020
	Attribution-NonCommercial-NoDerivatives 4.0 International
	You can use it for education only but can't make profits for any companies and individuals!
*/
package gomap

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	rw := make(map[string]interface{})

	// loop times
	var num = 1000
	randNum := 100000000
	rand.Seed(time.Now().Unix())
	// 1. new a map
	m := New()
	m = NewAVLMap()
	m.SetHash()
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		//fmt.Println("add key:", key)
		// 2. put key pairs
		xx := key + fmt.Sprintf("_%v", rand.Int63n(int64(randNum)))
		m.Put(key, xx)
		rw[key] = xx
		if m.Check() {
			//fmt.Println("is a rb tree,len:", m.Len())
		} else {
			fmt.Println("add")
			return
			// check rb tree
		}
	}

	if m.Check() {
		fmt.Println("is a rb tree,len:", m.Len())
	}

	for k, v := range rw {
		vv, ok := m.Get(k)
		if !ok {
			fmt.Println("err")
			return
		}

		if vv != v {
			fmt.Println("1 err", vv, v)
			return
		}
	}

	// 8. delete many
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		//fmt.Println("delete key:", key)
		m.Delete(key)
		delete(rw, key)
		if m.Check() {
			//fmt.Println("is a rb tree,len:", m.Len())
		} else {
			return
			// check rb tree
		}
	}

	for k, v := range rw {
		vv, ok := m.Get(k)
		if !ok {
			fmt.Println("err")
			return
		}

		if vv != v {
			fmt.Println("err", vv, v)
			return
		}
	}

	// 3. can iterator
	//iterator := m.Iterator()
	//for iterator.HasNext() {
	//	k, v := iterator.Next()
	//	fmt.Printf("Iterator key:%s,value %v\n", k, v)
	//}

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
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		//fmt.Println("delete key:", key)
		xx := key + fmt.Sprintf("_%v", rand.Int63n(int64(randNum)))
		m.Put(key, xx)
		rw[key] = xx
		if m.Check() {
			//	//fmt.Println("is a rb tree,len:", m.Len())
		} else {
			// check rb tree
		}

		key = fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		m.Delete(key)
		delete(rw, key)
		if m.Check() {
			//	//fmt.Println("is a rb tree,len:", m.Len())
		} else {
			return
			// check rb tree
		}
	}

	for k, v := range rw {
		vv, ok := m.Get(k)
		if !ok {
			fmt.Println("err")
			return
		}

		if vv != v {
			fmt.Println("err", vv, v)
			return
		}
	}

	// 9. key list
	//fmt.Printf("keyList:%#v,len:%d\n", m.KeyList(), m.Len())
	//fmt.Printf("keySortList:%#v,len:%d\n", m.KeySortedList(), m.Len())

	// 10. check is a rb tree
	if m.Check() {
		fmt.Println("is a rb tree,len:", m.Len())
	}
}
