package data_struct

import "fmt"

type ListNode struct {
	Next  *ListNode
	Value interface{}
}

type LinkedList struct {
	Head   *ListNode
	length uint
}

func NewListNode(v interface{}) *ListNode {
	return &ListNode{nil, v}
}

func (this *ListNode) GetNext() *ListNode {
	return this.Next
}

func (this *ListNode) GetValue() interface{} {
	return this.Value
}

func NewLinkedList() *LinkedList {
	return &LinkedList{NewListNode(0), 0}
}

//在某个节点后面插入节点
func (this *LinkedList) InsertAfter(p *ListNode, v interface{}) bool {
	if nil == p {
		return false
	}
	newNode := NewListNode(v)
	oldNext := p.Next
	p.Next = newNode
	newNode.Next = oldNext
	this.length++
	return true
}

//在某个节点前面插入节点
func (this *LinkedList) InsertBefore(p *ListNode, v interface{}) bool {
	if nil == p || p == this.Head {
		return false
	}
	cur := this.Head.Next
	pre := this.Head
	for nil != cur {
		if cur == p {
			break
		}
		pre = cur
		cur = cur.Next
	}
	if nil == cur {
		return false
	}
	newNode := NewListNode(v)
	pre.Next = newNode
	newNode.Next = cur
	this.length++
	return true
}

//在链表头部插入节点
func (this *LinkedList) InsertToHead(v interface{}) bool {
	return this.InsertAfter(this.Head, v)
}

//在链表尾部插入节点
func (this *LinkedList) InsertToTail(v interface{}) bool {
	cur := this.Head
	for nil != cur.Next {
		cur = cur.Next
	}
	return this.InsertAfter(cur, v)
}

//通过索引查找节点
func (this *LinkedList) FindByIndex(index uint) *ListNode {
	if index >= this.length {
		return nil
	}
	cur := this.Head.Next
	var i uint = 0
	for ; i < index; i++ {
		cur = cur.Next
	}
	return cur
}

//删除传入的节点
func (this *LinkedList) DeleteNode(p *ListNode) bool {
	if nil == p {
		return false
	}
	cur := this.Head.Next
	pre := this.Head
	for nil != cur {
		if cur == p {
			break
		}
		pre = cur
		cur = cur.Next
	}
	if nil == cur {
		return false
	}
	pre.Next = p.Next
	p = nil
	this.length--
	return true
}

//打印链表
func (this *LinkedList) Print() {
	cur := this.Head.Next
	format := ""
	for nil != cur {
		format += fmt.Sprintf("%+v", cur.GetValue())
		cur = cur.Next
		if nil != cur {
			format += "->"
		}
	}
	fmt.Println(format)
}


/*
单链表反转
时间复杂度：O(N)
*/
func (this *LinkedList) Reverse() {
	if nil == this.Head || nil == this.Head.Next || nil == this.Head.Next.Next {
		return
	}

	var pre *ListNode = nil
	cur := this.Head.Next
	for nil != cur {
		tmp := cur.Next
		cur.Next = pre
		pre = cur
		cur = tmp
	}

	this.Head.Next = pre
}

/*
判断单链表是否有环
*/
func (this *LinkedList) HasCycle() bool {
	if nil != this.Head {
		slow := this.Head
		fast := this.Head
		for nil != fast && nil != fast.Next {
			slow = slow.Next
			fast = fast.Next.Next
			if slow == fast {
				return true
			}
		}
	}
	return false
}

/*
两个有序单链表合并
*/
func MergeSortedList(l1, l2 *LinkedList) *LinkedList {
	if nil == l1 || nil == l1.Head || nil == l1.Head.Next {
		return l2
	}
	if nil == l2 || nil == l2.Head || nil == l2.Head.Next {
		return l1
	}

	l := &LinkedList{Head: &ListNode{}}
	cur := l.Head
	curl1 := l1.Head.Next
	curl2 := l2.Head.Next
	for nil != curl1 && nil != curl2 {
		if curl1.Value.(int) > curl2.Value.(int) {
			cur.Next = curl2
			curl2 = curl2.Next
		} else {
			cur.Next = curl1
			curl1 = curl1.Next
		}
		cur = cur.Next
	}

	if nil != curl1 {
		cur.Next = curl1
	} else if nil != curl2 {
		cur.Next = curl2
	}

	return l
}

/*
删除倒数第N个节点
*/
func (this *LinkedList) DeleteBottomN(n int) {
	if n <= 0 || nil == this.Head || nil == this.Head.Next {
		return
	}

	fast := this.Head
	for i := 1; i <= n && fast != nil; i++ {
		fast = fast.Next
	}

	if nil == fast {
		return
	}

	slow := this.Head
	for nil != fast.Next {
		slow = slow.Next
		fast = fast.Next
	}
	slow.Next = slow.Next.Next
}

/*
获取中间节点
*/
func (this *LinkedList) FindMiddleNode() *ListNode {
	if nil == this.Head || nil == this.Head.Next {
		return nil
	}
	if nil == this.Head.Next.Next {
		return this.Head.Next
	}

	slow, fast := this.Head, this.Head
	for nil != fast && nil != fast.Next {
		slow = slow.Next
		fast = fast.Next.Next
	}
	return slow
}