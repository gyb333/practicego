package bench_test

import (
	"algorithm/sort"
	"math/rand"
	"testing"
	"time"
)
const length =10000000
var list []int

func InitData()  {
	// 以时间戳为种子生成随机数，保证每次运行数据不重复
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		list = append(list, int(r.Intn(length*5)))
	}
}
func init() {
	InitData()
}

func BenchmarkBubbleSort(b *testing.B) {
	var arr []int =make([]int,length/1000)
	copy(arr,list)
	sort.BubbleSort(arr)
}


func BenchmarkCocktailSort(b *testing.B) {
	var arr []int =make([]int,length/1000)
	copy(arr,list)
	sort.CocktailSort(arr)
}


func BenchmarkSelectionSort(b *testing.B) {
	var arr []int =make([]int,length/1000)
	copy(arr,list)
	sort.SelectionSort(arr)
}

func BenchmarkInsertionSort(b *testing.B) {
	var arr []int =make([]int,length/1000)
	copy(arr,list)
	sort.InsertionSort(arr)
}


func BenchmarkShellSort(b *testing.B) {
	var arr []int =make([]int,length/1000)
	copy(arr,list)
	sort.ShellSort(arr,3)
}

func BenchmarkQuickSort(b *testing.B) {
	var arr []int =make([]int,length)
	copy(arr,list)
	sort.QuickSort(arr)
}

func BenchmarkHeapSort(b *testing.B) {
	var arr []int =make([]int,length)
	copy(arr,list)
	sort.HeapSort(arr)
}

func BenchmarkMergeSort(b *testing.B) {
	var arr []int =make([]int,length)
	copy(arr,list)
	sort.MergeSort(arr,0,len(arr)-1)
}

func BenchmarkCountingSort(b *testing.B)  {
	var arr []int =make([]int,length)
	copy(arr,list)
	sort.CountingSort(arr)
}

func BenchmarkBucketSort(b *testing.B){
	var arr []int =make([]int,length)
	copy(arr,list)
	sort.BucketSort(arr)
}

func BenchmarkSliceIndexSort(b *testing.B) {
	var arr []int =make([]int,length)
	copy(arr,list)
	sort.SliceIndexSort(arr)
}

func BenchmarkRadixSort(b *testing.B)  {
	var arr []int =make([]int,length)
	copy(arr,list)
	sort.RadixSort(arr)
}