package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/kinix/proto/gchatpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var userName string

func main() {
	// Get the server address from env value
	addr := os.Getenv("GCHAT_SERVER")
	log.Printf("Connecting: %s\n", addr)

	// No SSL
	insecureOpt := grpc.WithTransportCredentials(insecure.NewCredentials())

	// Connect to the address
	clientConnection, err := grpc.Dial(addr, insecureOpt)
	if err != nil {
		log.Fatalf("could not connect: %s", err.Error())
		return
	}

	client := gchatpb.NewGChatClient(clientConnection)

	// The first message will be username
	fmt.Print("User Name: ")

	// Start the request
	stream, err := client.ReceiveAndSend(context.Background())
	if err != nil {
		log.Fatalf("could not open stream: %s", err.Error())
	}

	// For ending request
	close := make(chan bool)

	// For quiting from client
	quit := make(chan bool)

	go sendMessageFromStdin(stream, close, quit)
	go receiveMessageToStdin(stream, close)

	<-quit

	// Close connection
	clientConnection.Close()
}

func sendMessageFromStdin(stream gchatpb.GChat_ReceiveAndSendClient, close chan bool, quit chan bool) {
	var message string
	reader := bufio.NewReader(os.Stdin)

	for {
		// Even if close signal is caught, this selection will wait for stdin once more
		select {
		case <-close:
			quit <- true
			return
		default:
			// Read line from stdin (remove last character. It is EOL)
			message, _ = reader.ReadString('\n')
			stream.Send(&gchatpb.ChatMessage{
				Message: message[:len(message)-1],
			})
		}

	}

}

func receiveMessageToStdin(stream gchatpb.GChat_ReceiveAndSendClient, close chan bool) {
	var message *gchatpb.ChatMessage
	var err error

	for {
		message, err = stream.Recv()
		if err != nil {
			fmt.Println("Connection is closed.")
			close <- true
			return
		}

		// Print the message to stdin
		fmt.Printf("%s\n", message.GetMessage())
	}
}
