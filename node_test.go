package forest

import (
	"log"
	"testing"
	"time"
)

func TestNewJobNode(t *testing.T) {

	etcd := InitEtcd()

	jobNode, err := NewJobNode("192.168.10.35", etcd, ":8888", time.Second*5)
	if err != nil {
		t.Error(err)
	}

	log.Printf("the job Node:%#v", jobNode)

	go func() {

		time.Sleep(time.Second * 30)
		jobNode.Close()
	}()
	jobNode.Bootstrap()
}
func TestNewJobNode2(t *testing.T) {

	etcd := InitEtcd()

	jobNode, err := NewJobNode("192.168.10.36", etcd, ":8887", time.Second*5)
	if err != nil {
		t.Error(err)
	}

	log.Printf("the job Node:%#v", jobNode)
	go func() {

		time.Sleep(time.Second * 30)
		jobNode.Close()
	}()
	jobNode.Bootstrap()
}
