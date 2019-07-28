package forest

import (
	"context"
	"github.com/robfig/cron"
	"go.etcd.io/etcd/clientv3"
	"time"
)

const (
	KeyCreateChangeEvent = iota
	KeyUpdateChangeEvent
	KeyDeleteChangeEvent
)

const (
	JobCreateChangeEvent = iota
	JobUpdateChangeEvent
	JobDeleteChangeEvent
)

const (
	JobRunningStatus = iota + 1
	JobStopStatus
)

const (
	NodeFollowerState = iota
	NodeLeaderState
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
	Watcher    clientv3.Watcher
}

type TxResponse struct {
	Success bool
	LeaseID clientv3.LeaseID
	Lease   clientv3.Lease
	Key     string
	Value   string
}

// job
type JobConf struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Group  string `json:"group"`
	Cron   string `json:"cron"`
	Status int    `json:"status"`
	Target string `json:"target"`
	Params string `json:"params"`
	Mobile string `json:"mobile"`
	Remark string `json:"remark"`
}

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type GroupConf struct {
	Name   string `json:"name"`
	Remark string `json:"remark"`
}

type JobChangeEvent struct {
	Type int
	Conf *JobConf
}

type SchedulePlan struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Group      string `json:"group"`
	Cron       string `json:"cron"`
	Status     int    `json:"status"`
	Target     string `json:"target"`
	Params     string `json:"params"`
	Mobile     string `json:"mobile"`
	Remark     string `json:"remark"`
	schedule   cron.Schedule
	NextTime   time.Time `json:"nextTime"`
	BeforeTime time.Time `json:"beforeTime"`
}

type JobSnapshot struct {
	Id         string `json:"id"`
	JobId      string `json:"jobId"`
	Name       string `json:"name"`
	Ip         string `json:"ip"`
	Group      string `json:"group"`
	Cron       string `json:"cron"`
	Target     string `json:"target"`
	Params     string `json:"params"`
	Mobile     string `json:"mobile"`
	Remark     string `json:"remark"`
	CreateTime string `json:"createTime"`
}

type QueryClientParam struct {
	Group string `json:"group"`
}

type JobClient struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Group string `json:"group"`
}
type QuerySnapshotParam struct {
	Group string `json:"group"`
	Id    string `json:"id"`
	Ip    string `json:"ip"`
}

// node
type Node struct {
	Name  string `json:"name"`
	State int    `json:"state"`
}
