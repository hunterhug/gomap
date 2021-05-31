package gomap

import (
	"fmt"
	"math/rand"
	"testing"
)

var (
	randNum = 987654
)

// go test -run="bench_test.go" -test.bench=".*" -test.benchmem=1 -count=3
// -test.benchmem : 是否在性能测试的时候输出内存情况
// 循环次数， 平均每次执行时间
func BenchmarkGolangMapPut(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := make(map[string]interface{}, 0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		xx := key + fmt.Sprintf("_%v", rand.Int63n(int64(randNum)))
		m[key] = xx
	}
}

func BenchmarkRBTMapPut(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewMap()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		xx := key + fmt.Sprintf("_%v", rand.Int63n(int64(randNum)))
		m.Put(key, xx)
	}
}

func BenchmarkAVLMapPut(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewAVLMap()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		xx := key + fmt.Sprintf("_%v", rand.Int63n(int64(randNum)))
		m.Put(key, xx)
	}
}

func BenchmarkAVLRecursionMapPut(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewAVLRecursionMap()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		xx := key + fmt.Sprintf("_%v", rand.Int63n(int64(randNum)))
		m.Put(key, xx)
	}
}

func BenchmarkGolangMapDelete(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := make(map[string]interface{}, 0)
	for i := 0; i < randNum; i++ {
		key := fmt.Sprintf("%d", i)
		xx := key + fmt.Sprintf("_%v", i)
		m[key] = xx
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		delete(m, key)
	}
}

func BenchmarkRBTMapDelete(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewMap()
	for i := 0; i < randNum; i++ {
		key := fmt.Sprintf("%d", i)
		xx := key + fmt.Sprintf("_%v", i)
		m.Put(key, xx)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		m.Delete(key)
	}
}

func BenchmarkAVLMapDelete(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewAVLMap()
	for i := 0; i < randNum; i++ {
		key := fmt.Sprintf("%d", i)
		xx := key + fmt.Sprintf("_%v", i)
		m.Put(key, xx)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		m.Delete(key)
	}
}

func BenchmarkAVLRecursionMapDelete(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewAVLRecursionMap()
	for i := 0; i < randNum; i++ {
		key := fmt.Sprintf("%d", i)
		xx := key + fmt.Sprintf("_%v", i)
		m.Put(key, xx)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		m.Delete(key)
	}
}

func BenchmarkGolangMapGet(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := make(map[string]interface{}, 0)
	for i := 0; i < randNum; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		xx := key + fmt.Sprintf("_%v", rand.Int63n(int64(randNum)))
		m[key] = xx
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", -2) // can not fetch forever
		_ = m[key]
	}
}

func BenchmarkRBTMapGet(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewMap()
	for i := 0; i < randNum; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		xx := key + fmt.Sprintf("_%v", rand.Int63n(int64(randNum)))
		m.Put(key, xx)
	}

	//b.Logf("rb tree height:%d", m.Height())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", -2) // can not fetch forever
		_, _ = m.Get(key)
	}
}

func BenchmarkAVLMapGet(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewAVLMap()
	for i := 0; i < randNum; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		xx := key + fmt.Sprintf("_%v", rand.Int63n(int64(randNum)))
		m.Put(key, xx)
	}

	//b.Logf("avl tree height:%d", m.Height())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", -2) // can not fetch forever
		_, _ = m.Get(key)
	}
}

func BenchmarkAVLRecursionMapGet(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewAVLRecursionMap()
	for i := 0; i < randNum; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		xx := key + fmt.Sprintf("_%v", rand.Int63n(int64(randNum)))
		m.Put(key, xx)
	}

	//b.Logf("avl recur tree height:%d", m.Height())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", -2) // can not fetch forever
		_, _ = m.Get(key)
	}
}

func BenchmarkGolangMapRandom(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := make(map[string]interface{}, 0)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		xx := key + fmt.Sprintf("_%v", rand.Int63n(int64(randNum)))
		m[key] = xx

		key = fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		delete(m, key)
	}
}

func BenchmarkRBTMapRandom(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewMap()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		xx := key + fmt.Sprintf("_%v", rand.Int63n(int64(randNum)))
		m.Put(key, xx)

		key = fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		m.Delete(key)
	}
}

func BenchmarkAVLMapRandom(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewAVLMap()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		xx := key + fmt.Sprintf("_%v", rand.Int63n(int64(randNum)))
		m.Put(key, xx)

		key = fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		m.Delete(key)
	}
}

func BenchmarkAVLRecursionMapRandom(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewAVLRecursionMap()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		xx := key + fmt.Sprintf("_%v", rand.Int63n(int64(randNum)))
		m.Put(key, xx)

		key = fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		m.Delete(key)
	}
}
