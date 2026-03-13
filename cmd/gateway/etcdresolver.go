package main

import (
	"context"
	"strings"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

const etcdSchema = "etcd"

type etcdResolverBuilder struct {
	cli *clientv3.Client
}

func (e *etcdResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {
	ctx, cancel := context.WithCancel(context.Background())
	r := &etcdResolver{
		client:  e.cli,
		target:  target,
		cc:      cc,
		ctx:     ctx,
		cancel:  cancel,
		addrs:   make(map[string][]resolver.Address),
		watchCh: make(chan []resolver.Address, 1),
	}
	r.start()
	return r, nil
}

func (e *etcdResolverBuilder) Scheme() string {
	return etcdSchema
}

func registerEtcdResolverBuilder(cli *clientv3.Client) *etcdResolverBuilder {
	builder := &etcdResolverBuilder{
		cli: cli,
	}

	resolver.Register(builder)
	return builder
}

type etcdResolver struct {
	client  *clientv3.Client
	target  resolver.Target
	cc      resolver.ClientConn
	ctx     context.Context
	cancel  context.CancelFunc
	addrs   map[string][]resolver.Address
	watchCh chan []resolver.Address
	mu      sync.RWMutex
}

func (r *etcdResolver) ResolveNow(options resolver.ResolveNowOptions) {

}

func (r *etcdResolver) Close() {
	r.cancel()
}

func (r *etcdResolver) start() {
	serviceName := strings.TrimPrefix(r.target.URL.Path, "/")
	go r.watch(serviceName)
}

// serviceName: wallet-server, /wallet/server/
func (r *etcdResolver) watch(serviceName string) {
	splits := strings.Split(serviceName, "-")
	serviceName = strings.Join(splits, "/")
	serviceName = "/" + serviceName + "/"

	resp, err := r.client.Get(r.ctx, serviceName, clientv3.WithPrefix())
	if err == nil {
		addrs := make([]resolver.Address, 0, len(resp.Kvs))
		for _, kv := range resp.Kvs {
			addr := resolver.Address{Addr: string(kv.Value)}
			addrs = append(addrs, addr)
		}
		r.updateState(addrs)
	}

	watchCh := r.client.Watch(r.ctx, serviceName, clientv3.WithPrefix())
	for {
		select {
		case <-r.ctx.Done():
			return
		case resp := <-watchCh:
			for _, ev := range resp.Events {
				switch ev.Type {
				case clientv3.EventTypePut:
					r.updateAddress(ev)
				case clientv3.EventTypeDelete:
					r.updateAddress(ev)
				}
			}
		}
	}
}

func (r *etcdResolver) updateAddress(ev *clientv3.Event) {

}

func (r *etcdResolver) updateState(addrs []resolver.Address) {
	_ = r.cc.UpdateState(resolver.State{Addresses: addrs})
}
