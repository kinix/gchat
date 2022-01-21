package main

import (
	"log"
	"net"
	"os"

	"github.com/kinix/proto/gchatpb"
	"google.golang.org/grpc"
)

type server struct {
	gchatpb.UnimplementedGChatServer
	msgChannel msgChannelList
}

func main() {
	// Get the server address from env value
	addr := os.Getenv("GCHAT_SERVER")
	log.Printf("Server is starting: %s\n", addr)

	// Listen the address
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %s", err.Error())
	}

	// Create server object
	s := server{}

	// Create channel array (each channel will be used by one user in the chat)
	s.msgChannel = msgChannelList{
		list: []*chan string{},
	}

	gRPCserver := grpc.NewServer()

	// Register server
	gchatpb.RegisterGChatServer(gRPCserver, &s)

	// Start to serve
	if err := gRPCserver.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
