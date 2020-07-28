package sort

/*
原理
冒泡算法(bubble sort) 是一种很简单的交换排序。每轮都从第一个元素开始，依次将较大值向后交换一位，直至整个队列有序。
复杂度
和其他低效排序算法一样，平均时间复杂度是O(n^2)。最好的情况就是原队列就是排列好的数组，时间复杂度就是O(n)。空间复杂度为O(1)，用于交换。
按顺序通过比较排序的算法都是稳定的，冒泡排序也是这样。
 */
func BubbleSort(arr []int )  {
	for i:=1;i<len(arr);i++{
		for j:=0;j<len(arr)-i;j++{
			if arr[j]>arr[j+1]{
				arr[j],arr[j+1]=arr[j+1],arr[j]
			}
		}
	}
}

/*
原理
鸡尾酒排序(Cocktail Sort)是冒泡排序的一种优化算法。原本的冒泡排序只能在一轮中挑出一个值移动到最后，而鸡尾酒则可以在一轮里挑最大的移到最后，再挑最小的移到最前面。
实际上就是先正向进行一轮普通的冒泡排序，然后再逆向进行一轮反向冒泡，每轮冒泡都缩小一点范围。
复杂度
最好情况是正序排列的数列O(n)，最坏情况是逆序O(n^2)，平均还是O(n^2)。空间复杂度都是O(1)。
 */
func CocktailSort(list []int)  {
	var length=len(list)
	// 只需n/2轮的比较，因为每轮里都会吧最大值移到队尾，最小值移到队首
	for loop := 1; loop <= length/2; loop++ {
		sorted := false
		var j int
		// 先正向冒泡，把最大值移动到队尾
		for j = loop - 1; j < length-loop; j++ {
			if list[j] > list[j+1] {
				list[j], list[j+1] = list[j+1], list[j]
				sorted = true
			}
		}
		// 如果跑了一轮没有交换元素，说明已经排好序了
		if !sorted {
			break
		}
		// 再反向冒泡，把最小值移动到队首
		for ; j >= loop; j-- {
			if list[j] < list[j-1] {
				list[j], list[j-1] = list[j-1], list[j]
				sorted = true
			}
		}
		if !sorted {
			break
		}
	}
}
/*
原理
简单选择排序:直接从头开始一个一个去比，找出最小的放到最左边。再依次完成其他位的排序。

时间复杂度
比较次数固定为O(n^2)，数据交换次数是0~n-1次
因为会交换不同位置相同数值的数据，所以选择排序并不稳定
 */
func SelectionSort(arr []int )  {
	var minIndex int
	// 两次循环找到最小的排前面
	for i:=0;i<len(arr)-1;i++{
		minIndex=i
		for j:=i+1;j<len(arr);j++{
			if arr[j]<arr[minIndex]{
				minIndex=j
			}
		}
		// 当前就是最小值时就不交换了
		if minIndex == i {
			continue
		}
		// 顺序交换
		arr[i],arr[minIndex]=arr[minIndex],arr[i]
	}
}
/*
原理
直接插入排序:第一轮先从第二个元素开始，和第一个比较，如果较小就交换位置，本轮结束。
第二轮从第三个元素开始，先与第二个比较，如果较小就与第二个交换，交换后再于第一个比较。如此循环直至最后一个元素完成比较逻辑。

复杂度
最好的情况下，直接插入排序只需进行n-1次比较，0次的交换。平均下来时间复杂度为 O(n^2)。
由于是每个元素逐个与有序的队列进行比较，所以不会出现相同数值的元素在排序完成后交换位置。所以直接插入排序是种稳定的排序算法。

 */
func InsertionSort(arr []int)  {

	for i:=1;i<len(arr);i++{
		for j:=i;j>0 &&arr[j]<arr[j-1];j--{
			arr[j],arr[j-1]=arr[j-1],arr[j]
		}
		//fmt.Println(arr)
	}
}

/*
希尔排序:按某个增量值对数据进行分组，每组单独排序好后，再缩小这个增量，然后按新增量对数据分组后每个分组再各自排序。
最终增加缩小到1的时候，排序结束。所以希尔排序又叫缩小增量排序
复杂度:不同增量复杂度不同。n/2时平均的时间复杂度为O(n^2)。
相较直接插入排序，希尔排序减少了比较和交换的次数，在中小规模的排序中，性能表现较好。
但随着数据量增大，希尔排序与其他更好的排序算法（快排、堆排、并归等）仍有较大差距。
 */
func ShellSort(arr []int,gap int)  {
	length:=len(arr)
	//gap := 2
	step := length / gap

	for step >= 1 {
		// 这里按步长开始每个分组的排序
		for i := step; i < length; i++ {
			// 将按步长分组的子队列用直接插入排序算法进行排序
			insertionSortByStep(arr, step)
		}
		// 完成一轮后再次缩小增量
		step /= gap

		// 输出每轮缩小增量各组排序后的结果
		//fmt.Println(arr)
	}
}
func insertionSortByStep(tree []int, step int) {
	for i := step; i < len(tree); i++ {
		for j := i; j >= step && tree[j] < tree[j-step]; j -= step {
			tree[j], tree[j-step] = tree[j-step], tree[j]
		}
	}
}

/*
原理
快排原理其实比较简单，就是将原本很大的数组拆成小数组去解决问题。
要拆就得找个拆的位置。如果吧这个位置称为支点，那么快速排序问题就变成了不断的去找到拆分的支点元素位置。
通常找支点就是以某个元素为标准，通过交换元素位置把所有小于标准的元素都移到一侧，大于的移到另外一侧。移动元素的逻辑就是分别从最右侧元素向左找到比指定元素小的位置，再从最左侧开始向右找比指定元素大的位置。如果两个位置不相同就交换两个位置，在继续分表从两头相向寻找。找到合适的位置就是我们需要的支点。支点两边的元素再各自重复上面的操作，直到分拆出来的子数组只剩一个元素。分拆结束，顺序也就拍好了。
那么问题来了，以哪个元素为标准去比较呢？比如可以选第一个元素。
复杂度
理想情况下找到的支点可以把数组拆分成左右长度相近的子数组，此时时间复杂度为O(n*logn)
而最差情况则是每次找到的支点元素都在某一次，导致另一侧完全浪费，寻找支点的过程也浪费。这个时候用时会达到O(n^2)。
由于会打乱相同元素原有的顺序，所以快排也是一个不稳定排序。所以常用在普通类型数据的排序中。
 */
func QuickSort(arr []int){
	quickSort(arr,0,len(arr)-1)
}
func quickSort(list []int, start, end int)  {
	// 只剩一个元素时就返回了
	if start >= end {
		return
	}

	// 标记最左侧元素作为参考
	tmp := list[start]
	// 两个游标分别从两端相向移动，寻找合适的"支点"
	left := start
	right := end
	for left != right {
		// 右边的游标向左移动，直到找到比参考的元素值小的
		for list[right] >= tmp && left < right {
			right--
		}
		// 左侧游标向右移动，直到找到比参考元素值大的
		for list[left] <= tmp && left < right {
			left++
		}

		// 如果找到的两个游标位置不统一，就游标位置元素的值，并继续下一轮寻找
		// 此时交换的左右位置的值，右侧一定不大于左侧。可能相等但也会交换位置，所以才叫不稳定的排序算法
		if left < right {
			list[left], list[right] = list[right], list[left]
			//fmt.Println(list)
		}
	}

	// 这时的left位置已经是我们要找的支点了，交换位置
	list[start], list[left] = list[left], tmp

	// 按支点位置吧原数列分成两段，再各自逐步缩小范围排序
	quickSort(list, start, left-1)
	quickSort(list, left+1, end)
}

/*
堆排序就是利用大根堆（小根堆）的特性进行排序的。
从小到大排序一般用大根堆，从大到小一般用小根堆。

复杂度:平均o(n*logn)
由于初次构建大根堆时有较多次的排序，所以不适合对少量元素进行排序。由于相同数值的节点在比较过程中不能保证顺序，所以是种不稳定的排序方法。
 */

func HeapSort(arr []int){
	length:=len(arr)
	buildMaxHeap(arr)
	for i:=len(arr)-1;i>0;i--{
		arr[0],arr[i]=arr[i],arr[0]		//将最好一个元素跟根节点交换
		length--
		heapify(arr,length,0)
	}

}
func buildMaxHeap(arr []int){
	for i:=len(arr)/2;i>=0;i--{
		heapify(arr,len(arr),i)
	}
}
func heapify(arr []int,length,i int){
	 left:=2*i+1
	 right:=2*i+2
	 largest:=i;
	 if left<length&&arr[left]>arr[largest]{
	 	largest=left
	 }
	if right<length&&arr[right]>arr[largest]{
		largest=right
	}
	if largest!=i{
		arr[i],arr[largest]=arr[largest],arr[i]
		//fmt.Println(arr)
		heapify(arr,length,largest)
	}
}

func merge(arr []int, l int, mid int, r int) {
	temp := make([]int, r-l+1)
	for i := l; i <= r; i++ {
		temp[i-l] = arr[i]
	}
	left := l
	right := mid + 1

	for i := l; i <= r; i++ {
		if left > mid {
			arr[i] = temp[right-l]
			right++
		} else if right > r {
			arr[i] = temp[left-l]
			left++
		} else if temp[left - l] > temp[right - l] {
			arr[i] = temp[right - l]
			right++
		} else {
			arr[i] = temp[left - l]
			left++
		}
	}
}

func MergeSort(arr []int, l int, r int) {
	// 第二步优化，当数据规模足够小的时候，可以使用插入排序
	if r - l <= 15 {
		// 对 l,r 的数据执行插入排序
		for i := l + 1; i <= r; i++ {
			//temp := arr[i]
			//j := i
			//for ; j > 0 && temp < arr[j-1]; j-- {
			//	arr[j] = arr[j-1]
			//}
			//arr[j] = temp
			for j:=i;j>0 &&arr[j]<arr[j-1];j--{
				arr[j],arr[j-1]=arr[j-1],arr[j]
			}
		}
		return
	}
	if r>l{
		mid := (r + l) / 2
		MergeSort(arr, l, mid)
		MergeSort(arr, mid+1, r)

		// 第一步优化，左右两部分已排好序，只有当左边的最大值大于右边的最小值，才需要对这两部分进行merge操作
		if arr[mid] > arr[mid + 1] {
			merge(arr, l, mid, r)
		}
	}

}

// 计数排序
func CountingSort(data []int) {
	if len(data) <= 1 {
		return
	}
	min, max := countMaxMin(data)
	temp := make([]int, max+1)
	for i := 0; i < len(data); i++ {
		temp[data[i]]++
	}

	var index int
	for i := min; i < len(temp); i++ {
		for j := temp[i]; j > 0 ;j--{
			data[index] = i
			index++
		}
	}

}

func countMaxMin(data []int) (int, int) {
	min, max := data[0], data[0]
	for i := 1; i < len(data); i++ {
		if min > data[i] {
			min = data[i]
		}
		if max < data[i] {
			max = data[i]
		}
	}
	return min, max
}
/*
获取数组最大值
*/
func getMaxInArr(arr []int) int{
	max := arr[0]
	for i := 1; i < len(arr); i++ {
		if arr[i] > max{ max = arr[i]}
	}
	return max
}


func SliceIndexSort(arr []int) {
	t := make([]int,getMaxInArr(arr)+1)
	for _, val := range arr {
		t[val]++
	}
	res := make([]int, 0, len(arr))
	for index, val := range t {
		//循环把排序元素添加到新的数组中
		for ; val > 0; val-- {
			res = append(res, index)
		}
	}
	copy(arr,res)

}

/*
桶内排序
*/
func sortInBucket(bucket []int) {//此处实现插入排序方式，其实可以用任意其他排序方式
	length := len(bucket)
	if length == 1 {return}
	for i := 1; i < length; i++ {
		backup := bucket[i]
		j := i -1;
		//将选出的被排数比较后插入左边有序区
		for  j >= 0 && backup < bucket[j] {//注意j >= 0必须在前边，否则会数组越界
			bucket[j+1] = bucket[j]//移动有序数组
			j -- //反向移动下标
		}
		bucket[j + 1] = backup //插队插入移动后的空位
	}
}

/*
桶排序
*/
func BucketSort(arr []int)  {
	//桶数
	num := len(arr)
	//k（数组最大值）
	max := getMaxInArr(arr)
	//二维切片
	buckets := make([][]int, num)

	//分配入桶
	index := 0
	for i := 0; i < num; i++ {
		index = arr[i] * (num-1) /max//分配桶index = value * (n-1) /k
		buckets[index] = append(buckets[index], arr[i])
	}

	//桶内排序
	tmpPos := 0
	for i := 0; i < num; i++ {
		bucketLen := len(buckets[i])
		if bucketLen > 0{
			sortInBucket(buckets[i])
			copy(arr[tmpPos:], buckets[i])
			tmpPos += bucketLen
		}
	}
}

/*
基数排序(Radix Sort)是桶排序的扩展，它的基本思想是：将整数按位数切割成不同的数字，然后按每个位数分别比较。
具体做法是：将所有待比较数值统一为同样的数位长度，数位较短的数前面补零。
然后，从最低位开始，依次进行一次排序。这样从最低位排序一直到最高位排序完成以后, 数列就变成一个有序序列。
算法复杂度 时间复杂度为：O(K * N)
 */
func RadixSort(arr []int) []int {
	max := getMaxInArr(arr)
	// 数组中最大值决定了循环次数，101 循环三次
	for bit := 1; max/bit > 0; bit *= 10 {
		bitSort(arr, bit)
		//fmt.Println("[DEBUG bit]\t", bit)
		//fmt.Println("[DEBUG arr]\t", arr)
	}
	return arr
}
// 对指定的位进行排序
// bit 可取 1，10，100 等值
func bitSort(arr []int, bit int)  {
	n := len(arr)
	// 各个位的相同的数统计到 bitCounts[] 中
	bitCounts := make([]int, 10)
	for i := 0; i < n; i++ {
		num := (arr[i] / bit) % 10
		bitCounts[num]++
	}
	for i := 1; i < 10; i++ {
		bitCounts[i] += bitCounts[i-1]
	}

	tmp := make([]int, getMaxInArr(bitCounts))
	for i := n - 1; i >= 0; i-- {
		num := (arr[i] / bit) % 10
		tmp[bitCounts[num]-1] = arr[i]
		bitCounts[num]--
	}
	//for i := 0; i < n; i++ {
	//	arr[i] = tmp[i]
	//}
	copy(arr,tmp)
}
