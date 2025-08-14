package storage

import (
	"fmt"
	"sync"
	model "test-task-lo/internal/domain/models"
	"test-task-lo/internal/lib/random"
	st "test-task-lo/internal/storage"
)

type Storage struct {
	m  map[string]model.Task
	mu *sync.RWMutex
}

func New() *Storage {
	return &Storage{
		m:  make(map[string]model.Task, 32),
		mu: &sync.RWMutex{},
	}
}

func (s *Storage) SetTask(name string, desc string) (string, error) {
	id := random.NewRandomString(st.IDLength)
	task := model.Task{
		Name:   name,
		Desc:   desc,
		Status: "New",
		ID:     id,
	}
	s.mu.Lock()
	s.m[id] = task
	s.mu.Unlock()
	return id, nil
}

func (s *Storage) GetTask(id string) (model.Task, error) {
	const op = "storage.inmemorry-storage.GetTask"
	s.mu.RLock()
	task, ok := s.m[id]
	s.mu.RUnlock()
	if !ok {
		return model.Task{}, fmt.Errorf("%s, %w", op, st.ErrTaskNotFound)
	}
	return task, nil
}

func (s *Storage) GetTasks(status string) ([]model.Task, error) {
	s.mu.RLock()
	resp := MapGetValues(s.m, status)
	s.mu.RUnlock()
	return resp, nil
}

func MapGetValues(m map[string]model.Task, status string) []model.Task {
	values := make([]model.Task, 0, len(m))
	if status == "" {
		for _, v := range m {
			values = append(values, v)
		}
	} else {
		for _, v := range m {
			if v.Status == status {
				values = append(values, v)
			}

		}
	}

	return values
}
