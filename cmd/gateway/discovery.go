package main

import (
	"log"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var etcdAddr = os.Getenv("ETCD_ADDR")

func init() {
	if etcdAddr == "" {
		etcdAddr = "172.16.10.45:2379"
	}
}

type discovery struct {
	cli *clientv3.Client
}

func newDiscovery() *discovery {
	d := new(discovery)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdAddr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	d.cli = cli
	registerEtcdResolverBuilder(cli)
	return d
}

func (d *discovery) newGrpcConn() (*grpc.ClientConn, error) {
	target := "etcd:///wallet-server"
	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (d *discovery) stop() {
	_ = d.cli.Close()
}
