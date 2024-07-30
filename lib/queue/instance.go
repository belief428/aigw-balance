package queue

import (
	"sync"
)

// IQueue 队列信息
type IQueue interface {
	// Priority 优先级
	Priority() int
	// Delay 延迟时间
	Delay() int
	// Call 回调函数
	Call(args ...interface{})
}

type Node struct {
	Data IQueue
	Next *Node
}

type Instance struct {
	*Node
	size int
	lock *sync.RWMutex
}

func (this *Instance) LPush(data IQueue) {
	defer this.lock.Unlock()
	this.lock.Lock()

	node := &Node{Data: data}

	if this.IsEmpty() {
		this.Node = node
	} else {
		// 最开始节点
		cur := this.Node
		// 上一节点信息
		prev := new(Node)

		for cur != nil {
			if cur.Data.Priority() > node.Data.Priority() {
				prev = cur
				cur = cur.Next
				continue
			}
			node.Next = cur
			break
		}
		if prev.Data == nil {
			this.Node = node
		} else {
			prev.Next = node
		}
	}
	this.size++
}

func (this *Instance) RPush(data IQueue) {
	defer this.lock.Unlock()
	this.lock.Lock()

	node := &Node{Data: data}

	if this.IsEmpty() {
		this.Node = node
	} else {
		// 最开始节点
		cur := this.Node
		// 上一节点信息
		prev := new(Node)

		for cur != nil {
			if cur.Data.Priority() < node.Data.Priority() {
				break
			}
			prev = cur
			cur = cur.Next
		}
		node.Next = cur

		if prev.Data == nil {
			this.Node = node
		} else {
			prev.Next = node
		}
	}
	this.size++
}

func (this *Instance) LPop() IQueue {
	defer this.lock.Unlock()
	this.lock.Lock()

	if this.IsEmpty() {
		return nil
	}
	this.size--
	out := this.Node.Data
	this.Node = this.Node.Next
	return out
}

func (this *Instance) RProp() IQueue {
	defer this.lock.Unlock()
	this.lock.Lock()

	if this.IsEmpty() {
		return nil
	}
	out := this.Node

	cur := new(Node)

	for out.Next != nil {
		cur = out
		out = out.Next
	}
	if cur != nil {
		cur.Next = nil
	}
	this.size--
	return out.Data
}

func (this *Instance) IsEmpty() bool {
	return this.size <= 0
}

func (this *Instance) GetSize() int {
	return this.size
}

// NewInstance
// @Description:
// @return *Queue
func NewInstance() *Instance {
	return &Instance{
		size: 0,
		Node: nil,
		lock: new(sync.RWMutex),
	}
}
