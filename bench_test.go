package gomap

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// go test -run="bench_test.go" -test.bench=".*" -test.benchmem=1 -count=3
func BenchmarkBRTMapPut(b *testing.B) {
	b.StopTimer()

	randNum := 100000000
	rand.Seed(time.Now().Unix())

	m := NewMap()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		m.Put(key, key)
	}
}

func BenchmarkGolangMapPut(b *testing.B) {
	b.StopTimer()

	randNum := 100000000
	rand.Seed(time.Now().Unix())

	m := make(map[string]interface{}, 0)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		m[key] = key
	}
}

func BenchmarkBRTMapDelete(b *testing.B) {
	b.StopTimer()

	randNum := 100000000
	num := 100000
	rand.Seed(time.Now().Unix())

	m := NewMap()
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		xx := key + fmt.Sprintf("_%v", rand.Int63n(int64(randNum)))
		m.Put(key, xx)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		m.Delete(key)
	}
}

func BenchmarkGolangMapDelete(b *testing.B) {
	b.StopTimer()

	randNum := 100000000
	num := 100000
	rand.Seed(time.Now().Unix())

	m := make(map[string]interface{}, 0)
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		xx := key + fmt.Sprintf("_%v", rand.Int63n(int64(randNum)))
		m[key] = xx
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(randNum)))
		delete(m, key)
	}
}

func BenchmarkBRTMapRandom(b *testing.B) {
	b.StopTimer()

	randNum := 100000000

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

func BenchmarkGolangMapRandom(b *testing.B) {
	b.StopTimer()

	randNum := 100000000
	rand.Seed(time.Now().Unix())

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
