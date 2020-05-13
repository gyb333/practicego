package studymap

import (
	"testing"
	"strconv"
		)

const (
	MapSize     = 10000
	LoopCounter = 1000
)

/*
go test -bench . -count=5
go test -bench . -cpuprofile prof.cpu
go tool pprof [binary] prof.cpu对采样文件进行分析。
 */

func BenchmarkMap(b *testing.B) {
	// 初始化数据结构
	m := map[string]int{}
	keys := []string{}
	for i := 0; i < MapSize; i++ {
		keys = append(keys, strconv.Itoa(i))
		m[keys[i]] = i
	}
	// 重置测试计时器
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < LoopCounter; i++ {
			if _, ok := m[keys[i]]; ok {
			}
		}
	}
}
func BenchmarkMapWrite(b *testing.B) {
	// 初始化数据结构
	m := map[string]int{}
	keys := []string{}
	for i := 0; i < MapSize; i++ {
		keys = append(keys, strconv.Itoa(i))
		m[keys[i]] = i
	}
	// 重置测试计时器
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < LoopCounter; i++ {
			m[keys[i]] = i
		}
	}
}

func BenchmarkMapReadAndWrite(b *testing.B) {
	// 初始化数据结构
	m := map[string]int{}
	keys := []string{}
	for i := 0; i < MapSize; i++ {
		keys = append(keys, strconv.Itoa(i))
		m[keys[i]] = i
	}
	// 重置测试计时器
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < LoopCounter; i++ {
			if d, ok := m[keys[i]]; ok == true {
				m[keys[i]] = d + 1
			}
		}
	}

}

