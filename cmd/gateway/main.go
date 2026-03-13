package main

import (
	"context"
	"log"
	"net/http"

	pb "github.com/anylepride/wallet/proto/wallet"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	d := newDiscovery()
	defer d.stop()

	grpcConn, err := d.newGrpcConn()
	if err != nil {
		log.Fatal(err)
	}

	defer func(grpcConn *grpc.ClientConn) {
		_ = grpcConn.Close()
	}(grpcConn)

	mux := runtime.NewServeMux()
	err = pb.RegisterWalletServiceHandler(ctx, mux, grpcConn)
	if err != nil {
		log.Fatalf("failed to register wallet service handler: %v", err)
	}

	handler := corsMiddleWare(mux)

	log.Println("gRPC-gateway server listening on :8020")
	if err := http.ListenAndServe(":8020", handler); err != nil {
		log.Fatalf("failed to listen on port 8020, %v", err)
	}
}

func corsMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
