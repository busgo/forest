package forest

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type Etcd struct {
	endpoints []string
	client    *clientv3.Client
	kv        clientv3.KV

	timeout time.Duration
}

// create a etcd
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
		kv:        clientv3.NewKV(client),
		timeout:   timeout,
	}

	return
}

// get value  from a key
func (etcd *Etcd) Get(key string) (value []byte, err error) {

	var (
		getResponse *clientv3.GetResponse
	)
	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
	defer cancelFunc()

	if getResponse, err = etcd.kv.Get(ctx, key); err != nil {
		return
	}

	if len(getResponse.Kvs) == 0 {
		return
	}

	value = getResponse.Kvs[0].Value

	return

}

// get values from  prefixKey
func (etcd *Etcd) GetWithPrefixKey(prefixKey string) (keys [][]byte, values [][]byte, err error) {

	var (
		getResponse *clientv3.GetResponse
	)
	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
	defer cancelFunc()

	if getResponse, err = etcd.kv.Get(ctx, prefixKey, clientv3.WithPrefix()); err != nil {
		return
	}

	if len(getResponse.Kvs) == 0 {
		return
	}

	keys = make([][]byte, len(getResponse.Kvs))
	values = make([][]byte, len(getResponse.Kvs))

	for i := 0; i < len(getResponse.Kvs); i++ {
		keys = append(keys, getResponse.Kvs[i].Key)
		values = append(values, getResponse.Kvs[i].Value)
	}

	return

}

// put a key
func (etcd *Etcd) Put(key, value string) (err error) {

	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
	defer cancelFunc()

	if _, err = etcd.kv.Put(ctx, key, value); err != nil {
		return
	}

	return
}

// put a key not exist
func (etcd *Etcd) PutNotExist(key, value string) (success bool, oldValue []byte, err error) {

	var (
		txnResponse *clientv3.TxnResponse
	)
	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
	defer cancelFunc()

	txn := etcd.client.Txn(ctx)

	txnResponse, err = txn.If(clientv3.Compare(clientv3.Version(key), "=", 0)).

		Then(clientv3.OpPut(key, value)).
		Else(clientv3.OpGet(key)).
		Commit()

	if err != nil {
		return
	}

	if txnResponse.Succeeded {
		success = true
	} else {
		oldValue = make([]byte, 0)
		oldValue = txnResponse.Responses[0].GetResponseRange().Kvs[0].Value
	}

	return
}

func (etcd *Etcd) Update(key, value, oldValue string) (success bool, err error) {

	var (
		txnResponse *clientv3.TxnResponse
	)

	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
	defer cancelFunc()

	txn := etcd.client.Txn(ctx)

	txnResponse, err = txn.If(clientv3.Compare(clientv3.Value(key), "=", oldValue)).
		Then(clientv3.OpPut(key, value)).
		Commit()

	if err != nil {
		return
	}

	if txnResponse.Succeeded {
		success = true
	}

	return
}
