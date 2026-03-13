package main

import (
	"context"
	"log"

	pb "github.com/anylepride/wallet/proto/wallet"
	"github.com/anylepride/wallet/wallet"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) GenerateWallet(ctx context.Context, req *pb.CreateWalletReq) (*pb.CreateWalletResponse, error) {
	walletId := uuid.New().String()
	w := wallet.NewWallet(walletId)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.wallets[walletId] = w
	return &pb.CreateWalletResponse{WalletId: walletId}, nil
}

func (s *server) QueryWalletBalance(ctx context.Context, req *pb.QueryWalletReq) (*pb.QueryWalletResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	w, exist := s.wallets[req.WalletId]
	if !exist {
		return nil, status.Errorf(codes.NotFound, "wallet not found: %s", req.WalletId)
	}

	log.Printf("wallet: %v", w)
	return &pb.QueryWalletResponse{WalletId: w.GetWalletId(), Balance: w.GetBalance()}, nil
}

func (s *server) TransferWalletBalance(ctx context.Context, req *pb.TransferWalletReq) (*pb.TransferWalletResponse, error) {
	if req.SrcWalletId == "" || req.DestWalletId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid srcWalletId or destWalletId, %v %v", req.SrcWalletId, req.DestWalletId)
	}

	if req.Balance <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "bad balance, %v", req.Balance)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	srcWallet, exist := s.wallets[req.SrcWalletId]
	if !exist {
		return nil, status.Errorf(codes.NotFound, "srcWallet not found: %s", req.SrcWalletId)
	}

	destWallet, exist := s.wallets[req.DestWalletId]
	if !exist {
		return nil, status.Errorf(codes.NotFound, "destWallet not found: %s", req.DestWalletId)
	}

	srcBalance := srcWallet.GetBalance()
	if srcBalance < req.Balance {
		return nil, status.Errorf(codes.FailedPrecondition, "not enough balance, %v", req.Balance)
	}

	srcWallet.SubBalance(req.Balance)
	destWallet.AddBalance(req.Balance)
	return &pb.TransferWalletResponse{}, nil
}
