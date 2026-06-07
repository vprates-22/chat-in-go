package internal

import (
	"maps"
	"net"
	"sync"
)

type Topic struct {
	Title        string
	participants map[net.Conn]bool
	lock         sync.Mutex
}

func NewTopic(title string) *Topic {
	return &Topic{
		Title:        title,
		participants: make(map[net.Conn]bool),
		lock:         sync.Mutex{},
	}
}

func (t *Topic) AddParticipant(participant net.Conn, id string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.hasParticipant(participant) {
		return
	}

	println("User " + id + " joined topic " + t.Title)
	t.participants[participant] = true
}

func (t *Topic) RemoveParticipant(participant net.Conn, id string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if !t.hasParticipant(participant) {
		return
	}

	if id != "" {
		println("User " + id + " left topic " + t.Title)
	}
	delete(t.participants, participant)
}

func (t *Topic) BroadcastMessage(message string, sender net.Conn) {
	t.lock.Lock()
	participants := make(map[net.Conn]bool)
	maps.Copy(participants, t.participants)
	t.lock.Unlock()

	for participant := range participants {
		if _, err := participant.Write([]byte(message)); err != nil {
			t.RemoveParticipant(participant, "")
		}

	}
}

func (t *Topic) ParticipantCount() int {
	return len(t.participants)
}

func (t *Topic) hasParticipant(participant net.Conn) bool {
	_, ok := t.participants[participant]
	return ok
}
