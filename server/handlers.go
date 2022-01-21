package main

import (
	"fmt"

	"github.com/kinix/proto/gchatpb"
)

func (s *server) ReceiveAndSend(stream gchatpb.GChat_ReceiveAndSendServer) error {
	close := make(chan bool)

	msgChannel := make(chan string)

	go s.sendChatMessages(stream, &msgChannel, close)
	go s.receiveChatMessages(stream, close)

	<-close
	return nil
}

func (s *server) receiveChatMessages(stream gchatpb.GChat_ReceiveAndSendServer, close chan bool) {
	var userName string

	for {
		message, err := stream.Recv()
		if err != nil { // User is left
			if userName != "" {
				s.msgChannel.sendToEveryone(fmt.Sprintf("\n- %s is left.\n", userName))
			}
			// TODO: is it ok?
			close <- true
			return
		}

		if userName == "" {
			// The first message is username
			userName = message.GetMessage()
			s.msgChannel.sendToEveryone(fmt.Sprintf("\n- %s is joined.\n", userName))
		} else {
			// Send message (userName: message)
			s.msgChannel.sendToEveryone(fmt.Sprintf("%s: %s", userName, message.GetMessage()))
		}
	}
}

func (s *server) sendChatMessages(stream gchatpb.GChat_ReceiveAndSendServer, msgChannel *chan string, close chan bool) {
	// Add to server list in order to send the messages to the user
	s.msgChannel.add(msgChannel)

	var message string

	for {
		select {
		// Wait for new message
		case message = <-*msgChannel:
			// Send it to the user
			stream.Send(&gchatpb.ChatMessage{
				Message: message,
			})
		case <-close: // User is left
			// Remove it from the list
			s.msgChannel.remove(msgChannel)
			return
		}

	}
}
