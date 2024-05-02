package server

import (
	"context"
	"log"
	"net"

	authGrpc "git.ctisoftware.vn/back-end/base/src/grpc"

	"git.ctisoftware.vn/numerology/proto-lib/golang/license"
	"google.golang.org/grpc"
)

func ServeGrpc(ctx context.Context, addr string) (err error) {
	defer log.Println("GRPC server stopped", err)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	server := grpc.NewServer()
	license.RegisterLicenseServiceServer(server, &authGrpc.Server{})

	log.Printf("Listen and Serve License-Grpc-Service API at: %s\n", addr)
	return server.Serve(lis)
}
