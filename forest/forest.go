package main

import (
	"flag"
	"fmt"
	"github.com/busgo/forest"
	"github.com/prometheus/common/log"
	"strings"
	"time"
)

const (
	DefaultEndpoints   = "127.0.0.1:2379"
	DefaultHttpAddress = "%s:8888"
	DefaultDialTimeout = 5
)

func main() {

	ip := forest.GetLocalIpAddress()
	if ip == "" {
		log.Fatal("has no get the ip address")

	}

	fmt.Println(ip)

	endpoints := flag.String("etcd-endpoints", DefaultEndpoints, "etcd endpoints default 127.0.0.1:2379")
	httpAddress := flag.String("http-address", fmt.Sprintf(DefaultHttpAddress, ip), "etcd endpoints default :8888")
	etcdDialTime := flag.Int64("etcd-dailtimeout", DefaultDialTimeout, "provide the etcd dailtimeout default 5 second")
	help := flag.String("help", "", "forest help")
	flag.Parse()
	if *help != "" {
		flag.Usage()
		return
	}

	endpoint := strings.Split(*endpoints, ",")
	dialTime := time.Duration(*etcdDialTime) * time.Second

	etcd, err := forest.NewEtcd(endpoint, dialTime)
	if err != nil {
		log.Fatal(err)
	}

	node, err := forest.NewJobNode(ip, etcd, *httpAddress)
	if err != nil {

		log.Fatal(err)
	}

	node.Bootstrap()
}
