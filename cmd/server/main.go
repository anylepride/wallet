package main

import (
	"log"
	"net"

	pb "github.com/anylepride/wallet/proto/wallet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	grpcServer := grpc.NewServer()
	walletServer := NewServer()
	pb.RegisterWalletServiceServer(grpcServer, walletServer)

	reflection.Register(grpcServer)

	addr := getGrpcAddr()
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("grpc server listening on %v\n", addr)

	er := newEtcdRegister(addr)
	defer er.unregister(er.leaseId)

	if err := grpcServer.Serve(ln); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
