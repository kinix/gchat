package main

import "sync"

type msgChannelList struct {
	list  []*chan string
	mutex sync.RWMutex
}

// Send the message to all users
func (m *msgChannelList) sendToEveryone(msg string) {
	// TODO: Do not send to sender

	// Read lock for list
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, c := range m.list {
		*c <- msg
	}
}

func (m *msgChannelList) add(msgChannel *chan string) {
	// Lock for list
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.list = append(m.list, msgChannel)
}

func (m *msgChannelList) remove(msgChannel *chan string) {
	// Lock for list
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i, c := range m.list {
		if c == msgChannel {
			m.list = append(m.list[:i], m.list[i+1:]...)
		}
	}
}
