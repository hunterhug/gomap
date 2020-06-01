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

func BenchmarkRBTHashMapPut(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewMap().SetHash()
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

func BenchmarkAVLHashMapPut(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewAVLMap().SetHash()
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

func BenchmarkAVLRecursionHashMapPut(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewAVLRecursionMap().SetHash()
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

func BenchmarkRBTHashMapDelete(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewMap().SetHash()
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

func BenchmarkAVLHashMapDelete(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewAVLMap().SetHash()
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

func BenchmarkAVLRecursionHashMapDelete(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewAVLRecursionMap().SetHash()
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

func BenchmarkRBTHashMapGet(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewMap().SetHash()
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

func BenchmarkAVLHashMapGet(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewAVLMap().SetHash()
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

func BenchmarkAVLRecursionHashMapGet(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewAVLRecursionMap().SetHash()
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

func BenchmarkRBTHashMapRandom(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewMap().SetHash()
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

func BenchmarkAVLHashMapRandom(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewAVLMap().SetHash()
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

func BenchmarkAVLRecursionHashMapRandom(b *testing.B) {
	b.StopTimer()

	rand.Seed(int64(randNum))

	m := NewAVLRecursionMap().SetHash()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		xx := key + fmt.Sprintf("_%v", rand.Int63n(int64(randNum)))
		m.Put(key, xx)

		key = fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		m.Delete(key)
	}
}
