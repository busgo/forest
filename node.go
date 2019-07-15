package forest

import (
	"fmt"
	"log"
	"time"
)

const (
	JobNodePath = "/forest/server/node/"
	TTL         = 5
)

// job node
type JobNode struct {
	id           string
	registerPath string
	etcd         *Etcd
}

func NewJobNode(id string, etcd *Etcd) (node *JobNode, err error) {

	node = &JobNode{
		id:           id,
		registerPath: fmt.Sprintf("%s%s", JobNodePath, id),
		etcd:         etcd,
	}

	txResponse, err := node.registerJobNode()
	if err != nil {
		log.Fatalf("the job node:%s, fail register to :%s", node.id, node.registerPath)

	}

	if !txResponse.Success {
		log.Fatalf("the job node:%s, fail register to :%s,the job node id exist ", node.id, node.registerPath)
	}
	log.Printf("the job node:%s, success register to :%s", node.id, node.registerPath)
	go node.watchRegisterJobNode()

	return
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
		log.Printf("found the job node:%s register to path:%s has lose",node.id,node.registerPath)
		go node.loopRegisterJobNode()

	}
}

func (node *JobNode) registerJobNode() (txResponse *TxResponse, err error) {

	return node.etcd.TxKeepaliveWithTTL(node.registerPath, node.id, TTL)
}

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
			goto RETRY
		}
		log.Printf("the job node:%s,has already success register to :%s", node.id, node.registerPath)
	}

}
