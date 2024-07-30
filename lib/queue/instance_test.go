package queue

import (
	"fmt"
	"testing"
)

type Data struct {
	Key       string
	Priority_ int
}

func (this *Data) Priority() int {
	return this.Priority_
}

func (this *Data) Delay() int {
	return 0
}

func (this *Data) Call(args ...interface{}) {
	fmt.Println(this.Key)
}

var _queue = NewInstance()

func TestQueue_Put(t *testing.T) {
	//go func() {
	//	time.Sleep(1 * time.Second)
	//
	//	for {
	//		__queue := _queue.rprop()
	//
	//		if __queue == nil {
	//			time.Sleep(500 * time.Millisecond)
	//			continue
	//		}
	//		fmt.Println(__queue.Priority())
	//	}
	//}()

	//_bytes := make([]byte, 0)

	//_queue.RPush(&Data{
	//	Key:       "1",
	//	Priority_: 0,
	//})
	//_bytes, _ = json.Marshal(_queue)
	//t.Log(string(_bytes))
	//
	//_queue.RPush(&Data{
	//	Key:       "11",
	//	Priority_: 0,
	//})
	//_bytes, _ = json.Marshal(_queue)
	//t.Log(string(_bytes))
	//
	//_queue.LPush(&Data{
	//	Key:       "22",
	//	Priority_: 0,
	//})
	//_bytes, _ = json.Marshal(_queue)
	//t.Log(string(_bytes))
	////
	//RPush(&Data{
	//	Key:       "2",
	//	Priority_: 2,
	//})
	//_bytes, _ = json.Marshal(_queue)
	//t.Log(string(_bytes))
	//
	//RPush(&Data{
	//	Key:       "2",
	//	Priority_: 5,
	//})
	//_bytes, _ = json.Marshal(_queue)
	//t.Log(string(_bytes))

	//time.Sleep(10 * time.Second)

}
