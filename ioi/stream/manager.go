package stream

import "github.com/vmxy/go-ioi/ioi/util"

type Manager[T any] struct {
	conns *util.Map[string, T]
}

func NewManager[T any]() Manager[T] {
	m := util.NewMap[string, T]()
	manager := Manager[T]{
		conns: &m,
	}
	return manager
}
func (m *Manager[T]) Set(key string, session T) {
	m.conns.Set(key, session)
}
func (m *Manager[T]) Get(key string) (T, bool) {
	return m.conns.Get(key)
}
func (m *Manager[T]) Delete(key string) {
	m.conns.Delete(key)
}
func (m *Manager[T]) Size() int {
	return m.conns.Size()
}
func (m *Manager[T]) Rang(key string) (T, bool) {

	return m.conns.Get(key)
}
