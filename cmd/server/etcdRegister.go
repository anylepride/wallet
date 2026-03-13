package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"go.etcd.io/etcd/client/v3"
)

type etcdRegister struct {
	cli       *clientv3.Client
	serverKey string
	addr      string
	leaseId   clientv3.LeaseID
}

var etcdAddr = os.Getenv("ETCD_ADDR")

func init() {
	if etcdAddr == "" {
		etcdAddr = "172.16.10.45:2379"
	}
}

func newEtcdRegister(addr string) *etcdRegister {
	er := &etcdRegister{addr: addr, serverKey: uuid.New().String()}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdAddr},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}

	er.cli = cli
	er.register()
	return er
}

func (er *etcdRegister) register() {
	serviceName := fmt.Sprintf("/wallet/server/%s", er.serverKey)
	resp, err := er.cli.Grant(context.Background(), 10)
	if err != nil {
		log.Fatalf("etcd grant err, %v", err)
	}
	_, err = er.cli.Put(context.Background(), serviceName, er.addr, clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatalf("etcd register service err, %v", err)
	}
	log.Printf("etcd register service success, serviceName:%s addr:%v", serviceName, er.addr)
	// 保持心跳
	keepAliveChan, err := er.cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		log.Fatalf("etcd keep alive err, %v", err)
	}

	er.leaseId = resp.ID
	go er.keepAlive(keepAliveChan)
}

func (er *etcdRegister) keepAlive(c <-chan *clientv3.LeaseKeepAliveResponse) {
	for {
		select {
		case _, ok := <-c:
			if !ok {
				return
			}
		}
	}
}

func (er *etcdRegister) unregister(leaseId clientv3.LeaseID) {
	_, _ = er.cli.Revoke(context.Background(), leaseId)
	_ = er.cli.Close()
}
