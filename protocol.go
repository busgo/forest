package forest

import (
	"context"
	"go.etcd.io/etcd/clientv3"
)

const (
	KeyCreateChangeEvent = iota
	KeyUpdateChangeEvent
	KeyDeleteChangeEvent
)

// key 变化事件
type KeyChangeEvent struct {
	Type  int
	Key   string
	Value []byte
}

// 监听key 变化响应
type WatchKeyChangeResponse struct {
	Event      chan *KeyChangeEvent
	CancelFunc context.CancelFunc
	watcher    clientv3.Watcher
}

type TxResponse struct {
	Success bool
	LeaseID clientv3.LeaseID
	Lease   clientv3.Lease
	Key     string
	Value   string
}
