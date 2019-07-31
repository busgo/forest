package forest

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"time"
)

// fail over the job snapshot when the task client

type JobSnapshotFailOver struct {
	node                   *JobNode
	deleteClientEventChans chan *JobClientDeleteEvent
}

// new job snapshot fail over
func NewJobSnapshotFailOver(node *JobNode) (f *JobSnapshotFailOver) {

	f = &JobSnapshotFailOver{
		node:                   node,
		deleteClientEventChans: make(chan *JobClientDeleteEvent, 50),
	}

	f.loop()

	return
}

// loop
func (f *JobSnapshotFailOver) loop() {

	go func() {

		for ch := range f.deleteClientEventChans {
			f.handleJobClientDeleteEvent(ch)
		}
	}()
}

// handle job client delete event
func (f *JobSnapshotFailOver) handleJobClientDeleteEvent(event *JobClientDeleteEvent) {

	var (
		keys    [][]byte
		values  [][]byte
		err     error
		client  *Client
		success bool
	)

RETRY:
	prefixKey := fmt.Sprintf(JobClientSnapshotPath, event.Group.name, event.Client.name)
	if keys, values, err = f.node.etcd.GetWithPrefixKey(prefixKey); err != nil {
		log.Errorf("the fail client:%v for path:%s,error must retry", event.Client, prefixKey)
		time.Sleep(time.Second)
		goto RETRY
	}

	if len(keys) == 0 || len(values) == 0 {
		log.Warnf("the fail client:%v for path:%s is empty", event.Client, prefixKey)
		return
	}

	for pos := 0; pos < len(keys); pos++ {

		if client, err = event.Group.selectClient(); err != nil {
			log.Warnf("%v", err)
			continue
		}

		to := fmt.Sprintf(JobClientSnapshotPath, event.Group.name, client.name)

		from := string(keys[pos])
		value := string(values[pos])
		//  transfer the kv
		if success, _ = f.node.etcd.transfer(from, to, value); success {
			log.Infof("the fail client:%v for path:%s success transfer form %s to %s", event.Client, prefixKey, from, to)
		}

	}

}
