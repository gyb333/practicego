package studymap

import (
	"testing"
	"strconv"
	"sync"
	"runtime"
)

func BenchmarkRead_RWMap(b *testing.B) {
	m := NewRWMap()
	keys := []string{}
	for i := 0; i < MapSize; i++ {
		keys = append(keys, strconv.Itoa(i))
		m.Set(keys[i], i)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < LoopCounter; i++ {
			if _, ok := m.Get(keys[i]); ok {

			}
		}
	}
}

func BenchmarkWrite_RWMap(b *testing.B) {
	m := NewRWMap()
	keys := []string{}
	for i := 0; i < MapSize; i++ {
		keys = append(keys, strconv.Itoa(i))
		m.Set(keys[i], i)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < LoopCounter; i++ {
			m.Set(keys[i], i+1)
		}
	}
}

func BenchmarkReadAndWrite_RWMap(b *testing.B) {
	m := NewRWMap()
	keys := []string{}
	for i := 0; i < MapSize; i++ {
		keys = append(keys, strconv.Itoa(i))
		m.Set(keys[i], i)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < LoopCounter; i++ {
			if _, ok := m.Get(keys[i]); ok {
				m.Set(keys[i], i+1)
			}
		}
	}
}


func BenchmarkReadAndWriteWithMutilGoroutine_RWMap(b *testing.B) {
	m := NewRWMap()
	keys := []string{}
	for i := 0; i < MapSize; i++ {
		keys = append(keys, strconv.Itoa(i))
		m.Set(keys[i], i)
	}
	b.ResetTimer()
	wg := sync.WaitGroup{}
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < b.N; n++ {
				for i := 0; i < LoopCounter; i++ {
					if _, ok := m.Get(keys[i]); ok {
						m.Set(keys[i], i+1)
					}
				}
			}
		}()
	}
	wg.Wait()

}