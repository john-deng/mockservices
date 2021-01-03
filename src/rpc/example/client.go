//go:generate protoc --proto_path=./pkg/proto/user --go_out=plugins=grpc:./pkg/proto/user --go_opt=paths=source_relative ./pkg/proto/user/user.proto


package main
import (
	"context"
	"log"
	"strings"
	"time"

	"solarmesh.io/mockservices/src/grpc/example/pkg/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	log.Println("Client running ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := user.NewDeleteServiceClient(conn)

	request := &user.DeleteRequest{Uuid: "UUID-123"}

	// Anything linked to this variable will transmit request headers.
	md := metadata.New(map[string]string{"x-request-id": "req-123"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Anything linked to this variable will fetch response headers.
	var header metadata.MD

	response, err := client.Delete(ctx, request, grpc.Header(&header))
	if err != nil {
		log.Fatalln(err)
	}

	xrid := header["x-response-id"]
	if len(xrid) == 0 {
		log.Fatalln("missing 'x-response-id' header")
	}
	if strings.Trim(xrid[0], " ") == "" {
		log.Fatalln("empty 'x-response-id' header")
	}

	log.Println(response.GetCode())
	log.Println(xrid[0])
}