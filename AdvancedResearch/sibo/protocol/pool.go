package protocol

import "sync"

/*
在实际场景中每个message的结构都是一样的只是具体的内容值不同，
而且每个message的header是大小结构是相同的，因此我们可以将每个message的header缓存起来复用
 */
var msgPool = sync.Pool{
	New: func() interface{} {
		header := Header([8]byte{})
		header[0] = magicNumber

		return &Message{Header: &header}
	},
}

// GetPooledMsg gets a pooled message.
func GetPooledMsg() *Message {
	return msgPool.Get().(*Message)
}

// FreeMsg puts a msg into the pool.
func FreeMsg(msg *Message) {
	msg.Reset()
	msgPool.Put(msg)
}

var poolUint32Data = sync.Pool{
	New: func() interface{} {
		data := make([]byte, 4)
		return &data
	},
}
