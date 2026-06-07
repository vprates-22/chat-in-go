package internal

import (
	"net"
	"sync"
)

type Channel struct {
	Clients map[net.Conn]bool
	Topics  map[string]*Topic
	lock    sync.Mutex
}

type TopicInfo struct {
	Key          string
	Participants int
}

var mutexChannel = sync.Mutex{}
var singletonChannel *Channel

func NewChannel() *Channel {
	mutexChannel.Lock()
	defer mutexChannel.Unlock()

	if singletonChannel != nil {
		return singletonChannel
	}

	singletonChannel = &Channel{
		Clients: make(map[net.Conn]bool),
		Topics:  make(map[string]*Topic),
		lock:    sync.Mutex{},
	}

	return singletonChannel
}

func (c *Channel) AddClient(client net.Conn) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.Clients[client] = true
}

func (c *Channel) AddParticipant(key string, participant net.Conn, id string) {
	c.getOrCreateTopic(key).AddParticipant(participant, id)
}

func (c *Channel) RemoveParticipant(key string, participant net.Conn, id string) {
	c.getOrCreateTopic(key).RemoveParticipant(participant, id)

	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.Clients, participant)
}

func (c *Channel) ListTopics() []TopicInfo {
	c.lock.Lock()
	defer c.lock.Unlock()

	infos := make([]TopicInfo, 0, len(c.Topics))
	for k, v := range c.Topics {
		infos = append(infos, TopicInfo{
			Key:          k,
			Participants: v.ParticipantCount(),
		})
	}
	return infos
}

func (c *Channel) BroadcastMessage(key string, message string, sender net.Conn) {
	c.getOrCreateTopic(key).BroadcastMessage(message, sender)
}

func (c *Channel) Shutdown() {
	c.lock.Lock()
	defer c.lock.Unlock()

	for client := range c.Clients {
		client.Close()
	}
}

func (c *Channel) topicExists(key string) bool {
	_, ok := c.Topics[key]
	return ok
}

func (c *Channel) addTopic(title string) {
	newTopic := NewTopic(title)
	c.Topics[newTopic.Title] = newTopic
}

func (c *Channel) getOrCreateTopic(key string) *Topic {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.topicExists(key) {
		c.addTopic(key)
	}
	return c.Topics[key]
}
