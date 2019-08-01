package forest

import (
	"errors"
	"github.com/labstack/gommon/log"
	"strings"
)

const (
	JobConfPath = "/forest/server/conf/"
)

type JobManager struct {
	node *JobNode
}

func NewJobManager(node *JobNode) (manager *JobManager) {

	manager = &JobManager{
		node: node,
	}

	go manager.watchJobConfPath()

	return

}

func (manager *JobManager) watchJobConfPath() {

	keyChangeEventResponse := manager.node.etcd.WatchWithPrefixKey(JobConfPath)

	for ch := range keyChangeEventResponse.Event {

		manager.handleJobConfChangeEvent(ch)
	}
}

func (manager *JobManager) loopLoadJobConf() {

RETRY:
	var (
		keys   [][]byte
		values [][]byte
		err    error
	)
	if keys, values, err = manager.node.etcd.GetWithPrefixKey(JobConfPath); err != nil {

		goto RETRY
	}

	if len(keys) == 0 {
		return
	}

	for i := 0; i < len(keys); i++ {
		jobConf, err := UParkJobConf(values[i])
		if err != nil {
			log.Warnf("upark the job conf error:%#v", err)
			continue
		}
		manager.node.scheduler.pushJobChangeEvent(&JobChangeEvent{
			Type: JobCreateChangeEvent,
			Conf: jobConf,
		})

	}

}

func (manager *JobManager) handleJobConfChangeEvent(changeEvent *KeyChangeEvent) {

	switch changeEvent.Type {

	case KeyCreateChangeEvent:

		manager.handleJobCreateEvent(changeEvent.Value)
	case KeyUpdateChangeEvent:

		manager.handleJobUpdateEvent(changeEvent.Value)
	case KeyDeleteChangeEvent:

		manager.handleJobDeleteEvent(changeEvent.Key)

	}
}

func (manager *JobManager) handleJobCreateEvent(value []byte) {

	var (
		err     error
		jobConf *JobConf
	)
	if len(value) == 0 {

		return
	}

	if jobConf, err = UParkJobConf(value); err != nil {
		log.Printf("unpark the job conf err:%#v", err)
		return
	}

	manager.node.scheduler.pushJobChangeEvent(&JobChangeEvent{
		Type: JobCreateChangeEvent,
		Conf: jobConf,
	})

}

func (manager *JobManager) handleJobUpdateEvent(value []byte) {

	var (
		err     error
		jobConf *JobConf
	)
	if len(value) == 0 {

		return
	}

	if jobConf, err = UParkJobConf(value); err != nil {
		log.Printf("unpark the job conf err:%#v", err)
		return
	}

	manager.node.scheduler.pushJobChangeEvent(&JobChangeEvent{
		Type: JobUpdateChangeEvent,
		Conf: jobConf,
	})

}

// handle the job delete event
func (manager *JobManager) handleJobDeleteEvent(key string) {

	if key == "" {
		return
	}

	pos := strings.LastIndex(key, "/")
	if pos == -1 {
		return
	}

	id := key[pos+1:]

	jobConf := &JobConf{
		Id: id,
	}

	manager.node.scheduler.pushJobChangeEvent(&JobChangeEvent{
		Type: JobDeleteChangeEvent,
		Conf: jobConf,
	})

}

// add job conf
func (manager *JobManager) AddJob(jobConf *JobConf) (err error) {

	var (
		value   []byte
		v       []byte
		success bool
	)

	if value, err = manager.node.etcd.Get(GroupConfPath + jobConf.Group); err != nil {
		return
	}

	if len(value) == 0 {

		err = errors.New("任务集群不存在")
		return
	}

	jobConf.Id = GenerateSerialNo()
	jobConf.Version = 1

	if v, err = ParkJobConf(jobConf); err != nil {
		return
	}
	if success, _, err = manager.node.etcd.PutNotExist(JobConfPath+jobConf.Id, string(v)); err != nil {
		return
	}

	if !success {
		err = errors.New("创建失败,请重试！")
		return
	}
	return
}

// edit job conf
func (manager *JobManager) editJob(jobConf *JobConf) (err error) {

	var (
		value   []byte
		v       []byte
		success bool
		oldConf *JobConf
	)

	if value, err = manager.node.etcd.Get(GroupConfPath + jobConf.Group); err != nil {
		return
	}

	if len(value) == 0 {

		err = errors.New("任务集群不存在")
		return
	}

	if jobConf.Id == "" {
		err = errors.New("此记录任务配置记录不存在")
		return
	}



	if value, err = manager.node.etcd.Get(JobConfPath + jobConf.Id); err != nil {
		return
	}

	if len(value) == 0 {
		err = errors.New("此任务配置记录不存在")
		return
	}

	if oldConf, err = UParkJobConf([]byte(value)); err != nil {
		return
	}

	jobConf.Version =oldConf.Version+1
	if v, err = ParkJobConf(jobConf); err != nil {
		return
	}

	if success, err = manager.node.etcd.Update(JobConfPath+jobConf.Id, string(v), string(value)); err != nil {
		return
	}

	if !success {
		err = errors.New("修改失败,请重试！")
		return
	}
	return
}

// delete job conf
func (manager *JobManager) deleteJob(jobConf *JobConf) (err error) {

	var (
		value []byte
	)

	if jobConf.Id == "" {
		err = errors.New("此记录任务配置记录不存在")
		return
	}

	if value, err = manager.node.etcd.Get(JobConfPath + jobConf.Id); err != nil {
		return
	}

	if len(value) == 0 {
		err = errors.New("此任务配置记录不存在")
		return
	}
	err = manager.node.etcd.Delete(JobConfPath + jobConf.Id)

	return
}

// job list
func (manager *JobManager) jobList() (jobConfs []*JobConf, err error) {

	var (
		keys   [][]byte
		values [][]byte
	)
	if keys, values, err = manager.node.etcd.GetWithPrefixKey(JobConfPath); err != nil {
		return
	}

	if len(keys) == 0 {
		return
	}

	jobConfs = make([]*JobConf, 0)
	for i := 0; i < len(values); i++ {

		jobConf, err := UParkJobConf(values[i])
		if err != nil {
			log.Printf("upark the job conf errror:%#v", err)
			continue
		}

		jobConfs = append(jobConfs, jobConf)

	}

	return
}

// add group
func (manager *JobManager) addGroup(groupConf *GroupConf) (err error) {

	var (
		value   []byte
		success bool
	)
	if value, err = ParkGroupConf(groupConf); err != nil {
		return
	}

	if success, _, err = manager.node.etcd.PutNotExist(GroupConfPath+groupConf.Name, string(value)); err != nil {
		return
	}

	if !success {

		err = errors.New("此任务集群已存在")
	}

	return

}

// group list
func (manager *JobManager) groupList() (groupConfs []*GroupConf, err error) {

	var (
		keys   [][]byte
		values [][]byte
	)
	if keys, values, err = manager.node.etcd.GetWithPrefixKey(GroupConfPath); err != nil {
		return
	}

	if len(keys) == 0 {
		return
	}

	groupConfs = make([]*GroupConf, 0)
	for i := 0; i < len(values); i++ {

		groupConf, err := UParkGroupConf(values[i])
		if err != nil {
			log.Printf("upark the group conf errror:%#v", err)
			continue
		}

		groupConfs = append(groupConfs, groupConf)

	}

	return
}

// node list
func (manager *JobManager) nodeList() (nodes []string, err error) {

	var (
		keys   [][]byte
		values [][]byte
	)
	if keys, values, err = manager.node.etcd.GetWithPrefixKey(JobNodePath); err != nil {
		return
	}

	if len(keys) == 0 {
		return
	}

	nodes = make([]string, 0)
	for i := 0; i < len(values); i++ {
		nodes = append(nodes, string(values[i]))
	}

	return
}
