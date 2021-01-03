package main

import (
	"context"
	"log"
	"net"
	"strings"

	"solarmesh.io/mockservices/src/grpc/example/pkg/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type server struct {
	user.UnimplementedDeleteServiceServer
}

func main() {
	log.Println("Server running ...")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}

	srv := grpc.NewServer()
	user.RegisterDeleteServiceServer(srv, &server{})

	log.Fatalln(srv.Serve(lis))
}

func (s *server) Delete(ctx context.Context, request *user.DeleteRequest) (*user.DeleteResponse, error) {
	// Anything linked to this variable will fetch request headers.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "failed to get metadata")
	}
	xrid := md["x-request-id"]
	if len(xrid) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "missing 'x-request-id' header")
	}
	if strings.Trim(xrid[0], " ") == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty 'x-request-id' header")
	}

	log.Println(request.GetUuid())
	log.Println(xrid[0])

	// Anything linked to this variable will transmit response headers.
	header := metadata.New(map[string]string{"x-response-id": "res-123"})
	if err := grpc.SendHeader(ctx, header); err != nil {
		return nil, status.Errorf(codes.Internal, "unable to send 'x-response-id' header")
	}

	return &user.DeleteResponse{Code: 123}, nil
}