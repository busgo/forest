package forest

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/robfig/cron"
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

const (
	JobExecuteSnapshotDoingStatus   = 1
	JobExecuteSnapshotSuccessStatus = 2
	JobExecuteSnapshotUnkonwStatus  = 3
	JobExecuteSnapshotErrorStatus   = -1
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

type JobClientDeleteEvent struct {
	Client *Client
	Group  *Group
}

// job
type JobConf struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Group   string `json:"group"`
	Cron    string `json:"cron"`
	Status  int    `json:"status"`
	Target  string `json:"target"`
	Params  string `json:"params"`
	Mobile  string `json:"mobile"`
	Remark  string `json:"remark"`
	Version int    `json:"version"`
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
	Version    int       `json:"version"`
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

type JobExecuteSnapshot struct {
	Id         string `json:"id",xorm:"pk"`
	JobId      string `json:"jobId",xorm:"job_id"`
	Name       string `json:"name",xorm:"name"`
	Ip         string `json:"ip",xorm:"ip"`
	Group      string `json:"group",xorm:"group"`
	Cron       string `json:"cron",xorm:"cron"`
	Target     string `json:"target",xorm:"target"`
	Params     string `json:"params",xorm:"params"`
	Mobile     string `json:"mobile",xorm:"mobile"`
	Remark     string `json:"remark",xorm:"remark"`
	CreateTime string `json:"createTime",xorm:"create_time"`
	StartTime  string `json:"startTime",xorm:"start_time"`
	FinishTime string `json:"finishTime",xorm:"finish_time"`
	Times      int    `json:"times",xorm:"times"`
	Status     int    `json:"status",xorm:"status"`
	Result     string `json:"result",xorm:"result"`
}

type QueryExecuteSnapshotParam struct {
	Group    string `json:"group"`
	Id       string `json:"id"`
	Ip       string `json:"ip"`
	JobId    string `json:"jobId"`
	Name     string `json:"name"`
	Status   int    `json:"status"`
	PageSize int    `json:"pageSize"`
	PageNo   int    `json:"pageNo"`
}

type PageResult struct {
	TotalPage  int         `json:"totalPage"`
	TotalCount int         `json:"totalCount"`
	List       interface{} `json:"list"`
}

//  manual execute job
type ManualExecuteJobParam struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
}
