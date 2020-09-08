package queue

import "container/list"

type MessageBuffer struct {
	buf *list.List
}

func NewMessageBuffer() *MessageBuffer{
	return &MessageBuffer{buf:list.New()}
}

func (mb *MessageBuffer) Size() int{
	return mb.buf.Len()
}

func (mb *MessageBuffer)FrontValue() interface{}{
	return mb.buf.Front().Value
}

func (mb *MessageBuffer)BackValue() interface{} {
	return mb.buf.Back().Value
}

func (mb *MessageBuffer) Front() *list.Element  {
	return mb.buf.Front()
}

func (mb *MessageBuffer) Back() *list.Element  {
	return mb.buf.Back()
}

func (mb *MessageBuffer) Add(e interface{})  {
	mb.buf.PushFront(e)
}

func (mb *MessageBuffer) IsEmpty() bool{
	return mb.buf.Len()==0
}

func (mb *MessageBuffer) Clear() error {
	mb.buf.Init()
	return nil
}