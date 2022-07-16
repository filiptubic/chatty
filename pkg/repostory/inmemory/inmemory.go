package inmemory

import (
	"chatty/pkg/model"
	"sync"
)

type InMemoryRepository struct {
	m       *sync.RWMutex
	history map[model.Channel][]model.Message
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		m:       &sync.RWMutex{},
		history: make(map[model.Channel][]model.Message),
	}
}

func (r *InMemoryRepository) SaveInHistory(ch model.Channel, m model.Message) {
	r.m.Lock()
	defer r.m.Unlock()
	r.history[ch] = append(r.history[ch], m)
}

func (r *InMemoryRepository) LoadHistory(ch model.Channel) []model.Message {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.history[ch]
}
