package forest

import (
	"github.com/labstack/gommon/log"
	"sync"
	"time"
)

// collection job execute status

const (
	JobExecuteStatusCollectionPath = "/forest/client/execute/snapshot/"
)

type JobCollection struct {
	node *JobNode
	lk   *sync.RWMutex
}

func NewJobCollection(node *JobNode) (c *JobCollection) {

	c = &JobCollection{
		node: node,
		lk:   &sync.RWMutex{},
	}

	c.watch()
	go c.loop()
	return
}

// watch
func (c *JobCollection) watch() {

	keyChangeEventResponse := c.node.etcd.WatchWithPrefixKey(JobExecuteStatusCollectionPath)
	log.Printf("the job collection success watch for path:%s ", JobExecuteStatusCollectionPath)
	go func() {

		for event := range keyChangeEventResponse.Event {

			c.handleJobExecuteStatusCollectionEvent(event)
		}

	}()
}

// handle the job execute status
func (c *JobCollection) handleJobExecuteStatusCollectionEvent(event *KeyChangeEvent) {

	if c.node.state == NodeFollowerState {
		return
	}

	switch event.Type {

	case KeyCreateChangeEvent:

		if len(event.Value) == 0 {
			return
		}

		executeSnapshot, err := UParkJobExecuteSnapshot(event.Value)

		if err != nil {
			log.Warnf("UParkJobExecuteSnapshot:%s fail,err:%#v ", event.Value, err)
			_ = c.node.etcd.Delete(event.Key)
			return
		}
		c.handleJobExecuteSnapshot(event.Key, executeSnapshot)

	case KeyUpdateChangeEvent:

		if len(event.Value) == 0 {
			return
		}

		executeSnapshot, err := UParkJobExecuteSnapshot(event.Value)

		if err != nil {
			log.Warnf("UParkJobExecuteSnapshot:%s fail,err:%#v ", event.Value, err)
			return
		}

		c.handleJobExecuteSnapshot(event.Key, executeSnapshot)

	case KeyDeleteChangeEvent:

	}
}

// handle job execute snapshot
func (c *JobCollection) handleJobExecuteSnapshot(path string, snapshot *JobExecuteSnapshot) {

	var (
		exist bool
		err   error
	)

	c.lk.Lock()
	defer c.lk.Unlock()
	if exist, err = c.checkExist(snapshot.Id); err != nil {
		log.Printf("check snapshot exist  error:%v", err)
		return
	}

	if exist {
		c.handleUpdateJobExecuteSnapshot(path, snapshot)
	} else {
		c.handleCreateJobExecuteSnapshot(path, snapshot)
	}

}

// handle create job execute snapshot
func (c *JobCollection) handleCreateJobExecuteSnapshot(path string, snapshot *JobExecuteSnapshot) {

	if snapshot.Status == JobExecuteSnapshotUnkonwStatus || snapshot.Status == JobExecuteSnapshotErrorStatus || snapshot.Status == JobExecuteSnapshotSuccessStatus {
		_ = c.node.etcd.Delete(path)
	}

	dateTime, err := ParseInLocation(snapshot.CreateTime)
	days := 0
	if err == nil {

		days = TimeSubDays(time.Now(), dateTime)

	}
	if snapshot.Status == JobExecuteSnapshotDoingStatus && days >= 3 {
		_ = c.node.etcd.Delete(path)
	}
	_, err = c.node.engine.Insert(snapshot)
	if err != nil {
		log.Printf("err:%#v", err)
	}
}

// handle update job execute snapshot
func (c *JobCollection) handleUpdateJobExecuteSnapshot(path string, snapshot *JobExecuteSnapshot) {

	if snapshot.Status == JobExecuteSnapshotUnkonwStatus || snapshot.Status == JobExecuteSnapshotErrorStatus || snapshot.Status == JobExecuteSnapshotSuccessStatus {
		_ = c.node.etcd.Delete(path)
	}

	dateTime, err := ParseInLocation(snapshot.CreateTime)
	days := 0
	if err == nil {

		days = TimeSubDays(time.Now(), dateTime)

	}
	if snapshot.Status == JobExecuteSnapshotDoingStatus && days >= 3 {
		_ = c.node.etcd.Delete(path)
	}

	_, _ = c.node.engine.Where("id=?", snapshot.Id).Cols("status", "finish_time", "times", "result").Update(snapshot)

}

// check the exist
func (c *JobCollection) checkExist(id string) (exist bool, err error) {

	var (
		snapshot *JobExecuteSnapshot
	)

	snapshot = new(JobExecuteSnapshot)

	if exist, err = c.node.engine.Where("id=?", id).Get(snapshot); err != nil {
		return
	}

	return

}

func (c *JobCollection) loop() {

	timer := time.NewTimer(10 * time.Minute)

	for {

		key := JobExecuteStatusCollectionPath
		select {
		case <-timer.C:

			timer.Reset(10 * time.Second)
			keys, values, err := c.node.etcd.GetWithPrefixKeyLimit(key, 1000)
			if err != nil {

				log.Warnf("collection loop err:%v ", err)
				continue
			}

			if len(keys) == 0 {

				continue

			}

			for pos := 0; pos < len(keys); pos++ {

				executeSnapshot, err := UParkJobExecuteSnapshot(values[pos])

				if err != nil {
					log.Warnf("UParkJobExecuteSnapshot:%s fail,err:%#v ", values[pos], err)
					_ = c.node.etcd.Delete(string(keys[pos]))
					continue
				}

				if executeSnapshot.Status == JobExecuteSnapshotSuccessStatus || executeSnapshot.Status == JobExecuteSnapshotErrorStatus || executeSnapshot.Status == JobExecuteSnapshotUnkonwStatus {
					path := string(keys[pos])
					c.handleJobExecuteSnapshot(path, executeSnapshot)
				}

			}



		}
	}

}
