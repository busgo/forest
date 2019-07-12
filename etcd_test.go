package forest

import (
	"log"
	"testing"
	"time"
)

func InitEtcd() *Etcd {

	etcd, err := NewEtcd([]string{"127.0.0.1:2379"}, time.Second*5)
	if err != nil {
		log.Fatal(err)
	}

	return etcd
}

func TestEtcd_Put(t *testing.T) {
	etcd := InitEtcd()

	err := etcd.Put("/echo", "echo-value")
	if err != nil {
		log.Fatal(err)
	}

	err = etcd.Put("/echo/one", "echo-value-one")
	if err != nil {
		log.Fatal(err)
	}

}
func TestEtcd_Get(t *testing.T) {

	etcd := InitEtcd()

	value, err := etcd.Get("/echo")
	if err != nil {

		log.Fatal(err)
	}

	log.Println("get a value:", string(value))
}

func TestEtcd_GetWithPrefixKey(t *testing.T) {
	etcd := InitEtcd()

	keys, values, err := etcd.GetWithPrefixKey("/echo")
	if err != nil {

		log.Fatal(err)
	}

	for i, key := range keys {

		log.Println("key:", string(key))
		log.Println("value:", string(values[i]))
	}

}

func TestEtcd_PutNotExist(t *testing.T) {

	etcd := InitEtcd()

	success, old, err := etcd.PutNotExist("/echo", "echo-value")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("success", success)

	log.Println("old", string(old))
}

func TestEtcd_Update(t *testing.T) {

	etcd := InitEtcd()

	value, err := etcd.Get("/echo")
	if err != nil {
		log.Fatal(err)
	}

	success, err := etcd.Update("/echo", "echo-2", string(value))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("success:", success)
}
