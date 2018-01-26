// Package inmem provides in-memory implementations of all the domain repositories.
package inmem

import (
	"errors"
	"nxy/testsocket/game/session"
	"sync"
)

type testRepository struct {
	mtx  sync.RWMutex
	test map[string]string
}

func (t *testRepository) Store(id string) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.test[id] = id
	return nil
}

func (t *testRepository) Remove(id string) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	delete(t.test, id)
	return nil
}

func (t *testRepository) Find(id string) (string, error) {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	if val, ok := t.test[id]; ok {
		return val, nil
	}
	return "", errors.New("get test mem error")
}

func (t *testRepository) FindAll() []string {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	c := make([]string, 0, len(t.test))
	for _, val := range t.test {
		c = append(c, val)
	}
	return c
}

// NewCargoRepository returns a new instance of a in-memory cargo repository.
func NewTestRepository() session.Repository {
	return &testRepository{
		test: make(map[string]string),
	}
}
