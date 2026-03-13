package main

import (
	"log"
	"net"
	"sync"

	pb "github.com/anylepride/wallet/proto/wallet"
	"github.com/anylepride/wallet/wallet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedWalletServiceServer
	mu      sync.RWMutex
	wallets map[string]*wallet.Wallet
}

func NewServer() *server {
	return &server{
		wallets: make(map[string]*wallet.Wallet),
	}
}

func main() {
	grpcServer := grpc.NewServer()
	walletServer := NewServer()
	pb.RegisterWalletServiceServer(grpcServer, walletServer)

	reflection.Register(grpcServer)

	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("grpc server listening on :8000")
	if err := grpcServer.Serve(ln); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
