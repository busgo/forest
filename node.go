package forest

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/gommon/log"
	"time"
)

const (
	JobNodePath      = "/forest/server/node/"
	JobNodeElectPath = "/forest/server/elect/leader"
	TTL              = 5
)

// job node
type JobNode struct {
	id           string
	registerPath string
	electPath    string
	etcd         *Etcd
	state        int
	apiAddress   string
	api          *JobAPi
	manager      *JobManager
	scheduler    *JobScheduler
	groupManager *JobGroupManager
	exec         *JobExecutor
	engine       *xorm.Engine
	collection   *JobCollection
	failOver     *JobSnapshotFailOver
	listeners    []NodeStateChangeListener
	close        chan bool
}

// node state change  listener
type NodeStateChangeListener interface {
	notify(int)
}

func NewJobNode(id string, etcd *Etcd, httpAddress, dbUrl string) (node *JobNode, err error) {

	engine, err := xorm.NewEngine("mysql", dbUrl)
	if err != nil {
		return
	}

	node = &JobNode{
		id:           id,
		registerPath: fmt.Sprintf("%s%s", JobNodePath, id),
		electPath:    JobNodeElectPath,
		etcd:         etcd,
		state:        NodeFollowerState,
		apiAddress:   httpAddress,
		close:        make(chan bool),
		engine:       engine,
		listeners:    make([]NodeStateChangeListener, 0),
	}

	node.failOver = NewJobSnapshotFailOver(node)

	node.collection = NewJobCollection(node)

	node.initNode()

	// create job executor
	node.exec = NewJobExecutor(node)
	// create  group manager
	node.groupManager = NewJobGroupManager(node)

	node.scheduler = NewJobScheduler(node)

	// create job manager
	node.manager = NewJobManager(node)

	// create a job http api
	node.api = NewJobAPi(node)

	return
}

func (node *JobNode) addListeners() {

	node.listeners = append(node.listeners, node.scheduler)

}

func (node *JobNode) changeState(state int) {

	node.state = state

	if len(node.listeners) == 0 {

		return
	}

	// notify all listener
	for _, listener := range node.listeners {

		listener.notify(state)
	}

}

// start register node
func (node *JobNode) initNode() {
	txResponse, err := node.registerJobNode()
	if err != nil {
		log.Fatalf("the job node:%s, fail register to :%s", node.id, node.registerPath)

	}
	if !txResponse.Success {
		log.Fatalf("the job node:%s, fail register to :%s,the job node id exist ", node.id, node.registerPath)
	}
	log.Printf("the job node:%s, success register to :%s", node.id, node.registerPath)
	node.watchRegisterJobNode()
	node.watchElectPath()
	go node.loopStartElect()

}

// bootstrap
func (node *JobNode) Bootstrap() {

	go node.groupManager.loopLoadGroups()
	go node.manager.loopLoadJobConf()

	<-node.close
}

func (node *JobNode) Close() {

	node.close <- true
}

// watch the register job node
func (node *JobNode) watchRegisterJobNode() {

	keyChangeEventResponse := node.etcd.Watch(node.registerPath)

	go func() {

		for ch := range keyChangeEventResponse.Event {
			node.handleRegisterJobNodeChangeEvent(ch)
		}
	}()

}

// handle the register job node change event
func (node *JobNode) handleRegisterJobNodeChangeEvent(changeEvent *KeyChangeEvent) {

	switch changeEvent.Type {

	case KeyCreateChangeEvent:

	case KeyUpdateChangeEvent:

	case KeyDeleteChangeEvent:
		log.Printf("found the job node:%s register to path:%s has lose", node.id, node.registerPath)
		go node.loopRegisterJobNode()

	}
}

func (node *JobNode) registerJobNode() (txResponse *TxResponse, err error) {

	return node.etcd.TxKeepaliveWithTTL(node.registerPath, node.id, TTL)
}

// loop register the job node
func (node *JobNode) loopRegisterJobNode() {

RETRY:

	var (
		txResponse *TxResponse
		err        error
	)
	if txResponse, err = node.registerJobNode(); err != nil {
		log.Printf("the job node:%s, fail register to :%s", node.id, node.registerPath)
		time.Sleep(time.Second)
		goto RETRY
	}

	if txResponse.Success {
		log.Printf("the job node:%s, success register to :%s", node.id, node.registerPath)
	} else {

		v := txResponse.Value
		if v != node.id {
			time.Sleep(time.Second)
			log.Fatalf("the job node:%s,the other job node :%s has already  register to :%s", node.id, v, node.registerPath)
		}
		log.Printf("the job node:%s,has already success register to :%s", node.id, node.registerPath)
	}

}

// elect the leader
func (node *JobNode) elect() (txResponse *TxResponse, err error) {

	return node.etcd.TxKeepaliveWithTTL(node.electPath, node.id, TTL)

}

// watch the job node elect path
func (node *JobNode) watchElectPath() {

	keyChangeEventResponse := node.etcd.Watch(node.electPath)

	go func() {

		for ch := range keyChangeEventResponse.Event {

			node.handleElectLeaderChangeEvent(ch)
		}
	}()

}

// handle the job node leader change event
func (node *JobNode) handleElectLeaderChangeEvent(changeEvent *KeyChangeEvent) {

	switch changeEvent.Type {

	case KeyDeleteChangeEvent:
		node.changeState(NodeFollowerState)
		node.loopStartElect()
	case KeyCreateChangeEvent:

	case KeyUpdateChangeEvent:

	}

}

// loop start elect
func (node *JobNode) loopStartElect() {

RETRY:
	var (
		txResponse *TxResponse
		err        error
	)
	if txResponse, err = node.elect(); err != nil {
		log.Printf("the job node:%s,elect  fail to :%s", node.id, node.electPath)
		time.Sleep(time.Second)
		goto RETRY
	}

	if txResponse.Success {
		node.changeState(NodeLeaderState)
		log.Printf("the job node:%s,elect  success to :%s", node.id, node.electPath)
	} else {
		v := txResponse.Value
		if v != node.id {
			log.Printf("the job node:%s,give up elect request because the other job node：%s elect to:%s", node.id, v, node.electPath)
			node.changeState(NodeFollowerState)
		} else {
			log.Printf("the job node:%s, has already elect  success to :%s", node.id, node.electPath)
			node.changeState(NodeLeaderState)
		}
	}

}
