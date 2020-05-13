package studymap

import (
	"testing"
	"sync"
	"strconv"
	"runtime"
)

func BenchmarkReadAndWrite_SyncMap(b *testing.B) {
	var m sync.Map

	keys := []string{}
	for i := 0; i < MapSize; i++ {
		keys = append(keys, strconv.Itoa(i))
		m.Store(keys[i], i)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < LoopCounter; i++ {
			if _, ok := m.Load(keys[i]); ok {
				m.Store(keys[i], i+1)
			}
		}
	}
}

func BenchmarkReadAndWriteWithMutilGoroutine_SyncMap(b *testing.B) {
	var m sync.Map

	keys := []string{}
	for i := 0; i < MapSize; i++ {
		keys = append(keys, strconv.Itoa(i))
		m.Store(keys[i], i)
	}

	b.ResetTimer()
	wg := sync.WaitGroup{}
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < b.N; n++ {
				for i := 0; i < LoopCounter; i++ {
					if _, ok := m.Load(keys[i]); ok {
						m.Store(keys[i], i+1)
					}
				}
			}
		}()
	}
	wg.Wait()
}

