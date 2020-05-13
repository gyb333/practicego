package sort_test

import (
	"../sort"
	"math/rand"
	"testing"
	"time"
)

const length =10000
var list []int

func InitData()  {
	// 以时间戳为种子生成随机数，保证每次运行数据不重复
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		list = append(list, int(r.Intn(length*10)))
	}
}
func init() {
	InitData()
}

func TestBubbleSort(t *testing.T) {
	var arr []int =make([]int,length)
	copy(arr,list)
	sort.BubbleSort(arr)
	//fmt.Println(arr)
}

func TestCocktailSort(t *testing.T) {
	var arr []int =make([]int,length)
	copy(arr,list)
	sort.CocktailSort(arr)
	//fmt.Println(arr)
}

func TestSelectionSort(t *testing.T) {
	var arr []int =make([]int,length)
	copy(arr,list)

	sort.SelectionSort(arr)
	//fmt.Println(arr)
}

func TestInsertionSort(t *testing.T) {
	var arr []int =make([]int,length)
	copy(arr,list)
	sort.InsertionSort(arr)
	//fmt.Println(arr)
}

func TestShellSort(t *testing.T) {
	var arr []int =make([]int,length/100)
	copy(arr,list)
	sort.ShellSort(arr,3)
	//fmt.Println(arr)
}

func TestQuickSort(t *testing.T) {
	var arr []int =make([]int,length)
	copy(arr,list)
	sort.QuickSort(arr)
	//fmt.Println(arr)
}

func TestHeapSort(t *testing.T) {
	var arr []int =make([]int,length)
	copy(arr,list)
	//arr :=[]int{8,4,12,7,35,9,22,41,2}
	sort.HeapSort(arr)
	//fmt.Println(arr)
}

func TestMergeSort(t *testing.T) {
	var arr []int =make([]int,length)
	copy(arr,list)
	//arr :=[]int{6,8,4,12,7,35,9,22,3,41,2,1,5,11,10,20}
	sort.MergeSort(arr,0,len(arr)-1)
	//fmt.Println(arr)
}

func TestCountingSort(t *testing.T) {
	var arr []int =make([]int,length)
	copy(arr,list)
	//arr :=[]int{6,8,4,12,7,35,9,22,3,41,2,1,5,11,10,20}
	sort.CountingSort(arr)
	//fmt.Println(arr)
}

func TestBucketSort(t *testing.T) {
	var arr []int =make([]int,length)
	copy(arr,list)
	//arr :=[]int{6,8,4,12,7,35,9,22,3,41,2,1,5,11,10,20}
	sort.BucketSort(arr)
	//fmt.Println(arr)
}

func TestSliceIndexSort(t *testing.T) {
	var arr []int =make([]int,length)
	copy(arr,list)
	//arr :=[]int{6,8,4,12,7,35,9,22,3,41,2,1,5,11,10,20}
	sort.SliceIndexSort(arr)
	//fmt.Println(arr)
}

func TestRadixSort(t *testing.T) {
	var arr []int =make([]int,length)
	copy(arr,list)
	//arr :=[]int{6,8,4,12,7,35,9,22,3,41,2,1,5,11,10,20}
	sort.RadixSort(arr)
	//fmt.Println(arr)
}
