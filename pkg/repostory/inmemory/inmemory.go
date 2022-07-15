package inmemory

import (
	repository "chatty/pkg/repostory"
	"sync"
)

type InMemoryRepository struct {
	m       *sync.RWMutex
	history map[repository.Channel][]repository.Message
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		m:       &sync.RWMutex{},
		history: make(map[repository.Channel][]repository.Message),
	}
}

func (r *InMemoryRepository) SaveInHistory(ch repository.Channel, m repository.Message) {
	r.m.Lock()
	defer r.m.Unlock()
	r.history[ch] = append(r.history[ch], m)
}

func (r *InMemoryRepository) LoadHistory(ch repository.Channel) []repository.Message {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.history[ch]
}
