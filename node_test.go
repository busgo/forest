package forest

import (
	"log"
	"testing"
	"time"
)

func TestNewJobNode(t *testing.T) {

	etcd := InitEtcd()

	jobNode, err := NewJobNode("192.168.10.35", etcd)
	if err != nil {
		t.Error(err)
	}

	log.Printf("the job Node:%#v", jobNode)

	for {
		time.Sleep(time.Second)
	}
}
