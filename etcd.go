package forest

import (
	"go.etcd.io/etcd/clientv3"
	"time"
)

type Etcd struct {
	endpoints []string
	client    *clientv3.Client
	kv clientv3.KV
}

func NewEtcd(endpoints []string, timeout time.Duration) (etcd *Etcd, err error) {

	var (
		client *clientv3.Client
	)

	conf := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	}
	if client, err = clientv3.New(conf); err != nil {
		return
	}

	etcd = &Etcd{

		endpoints: endpoints,
		client:    client,
		kv:clientv3.NewKV(client),
	}

	return
}
