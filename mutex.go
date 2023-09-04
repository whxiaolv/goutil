package goutil

import (
	"fmt"
	"log"
	"sync"
)

type CommonMap struct {
	sync.RWMutex
	m map[string]interface{}
}

type Tuple struct {
	Key string
	Val interface{}
}

func NewCommonMap(size int) *CommonMap {
	if size > 0 {
		return &CommonMap{m: make(map[string]interface{}, size)}
	} else {
		return &CommonMap{m: make(map[string]interface{})}
	}
}
func (s *CommonMap) GetValue(k string) (interface{}, bool) {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.m[k]
	return v, ok
}
func (s *CommonMap) Put(k string, v interface{}) {
	s.Lock()
	defer s.Unlock()
	s.m[k] = v
}
func (s *CommonMap) Iter() <-chan Tuple { // reduce memory
	ch := make(chan Tuple)
	go func() {
		s.RLock()
		for k, v := range s.m {
			ch <- Tuple{Key: k, Val: v}
		}
		close(ch)
		s.RUnlock()
	}()
	return ch
}
func (s *CommonMap) LockKey(k string) {
	s.Lock()
	if v, ok := s.m[k]; ok {
		s.m[k+"_lock_"] = true
		s.Unlock()
		switch v.(type) {
		case *sync.Mutex:
			v.(*sync.Mutex).Lock()
		default:
			log.Print(fmt.Sprintf("LockKey %s", k))
		}
	} else {
		s.m[k] = &sync.Mutex{}
		v = s.m[k]
		s.m[k+"_lock_"] = true
		v.(*sync.Mutex).Lock()
		s.Unlock()
	}
}
func (s *CommonMap) UnLockKey(k string) {
	s.Lock()
	if v, ok := s.m[k]; ok {
		switch v.(type) {
		case *sync.Mutex:
			v.(*sync.Mutex).Unlock()
		default:
			log.Print(fmt.Sprintf("UnLockKey %s", k))
		}
		delete(s.m, k+"_lock_") // memory leak
		delete(s.m, k)          // memory leak
	}
	s.Unlock()
}
func (s *CommonMap) IsLock(k string) bool {
	s.Lock()
	if v, ok := s.m[k+"_lock_"]; ok {
		s.Unlock()
		return v.(bool)
	}
	s.Unlock()
	return false
}
func (s *CommonMap) Keys() []string {
	s.Lock()
	keys := make([]string, len(s.m))
	defer s.Unlock()
	for k, _ := range s.m {
		keys = append(keys, k)
	}
	return keys
}
func (s *CommonMap) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = make(map[string]interface{})
}
func (s *CommonMap) Remove(key string) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.m[key]; ok {
		delete(s.m, key)
	}
}
func (s *CommonMap) AddUniq(key string) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.m[key]; !ok {
		s.m[key] = nil
	}
}
func (s *CommonMap) AddCount(key string, count int) {
	s.Lock()
	defer s.Unlock()
	if _v, ok := s.m[key]; ok {
		v := _v.(int)
		v = v + count
		s.m[key] = v
	} else {
		s.m[key] = 1
	}
}
func (s *CommonMap) AddCountInt64(key string, count int64) {
	s.Lock()
	defer s.Unlock()
	if _v, ok := s.m[key]; ok {
		v := _v.(int64)
		v = v + count
		s.m[key] = v
	} else {
		s.m[key] = count
	}
}
func (s *CommonMap) Add(key string) {
	s.Lock()
	defer s.Unlock()
	if _v, ok := s.m[key]; ok {
		v := _v.(int)
		v = v + 1
		s.m[key] = v
	} else {
		s.m[key] = 1
	}
}
func (s *CommonMap) Zero() {
	s.Lock()
	defer s.Unlock()
	for k := range s.m {
		s.m[k] = 0
	}
}
func (s *CommonMap) Contains(i ...interface{}) bool {
	s.Lock()
	defer s.Unlock()
	for _, val := range i {
		if _, ok := s.m[val.(string)]; !ok {
			return false
		}
	}
	return true
}
func (s *CommonMap) Get() map[string]interface{} {
	s.Lock()
	defer s.Unlock()
	m := make(map[string]interface{})
	for k, v := range s.m {
		m[k] = v
	}
	return m
}
