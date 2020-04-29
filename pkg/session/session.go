package session

import (
	"errors"
	"net"
	"sync"
)

type ConnID string

// Sessions keeps track of the TCP connections using a pool(map).
type Sessions struct {
	locker sync.RWMutex
	conns  map[ConnID]net.Conn
}


// Init creates a new object of Sessions struct and returns the pointer of it to be used by its receiver methods.
func Init() *Sessions{
	s := Sessions{
		locker: sync.RWMutex{},
		conns:  map[ConnID]net.Conn{},
	}

	return &s
}

// AddConn adds a new connection to the sessions map.
func (s *Sessions) AddConn(uniqueID ConnID, conn net.Conn) error {
	s.locker.RLock()
	_, ok := s.conns[uniqueID]
	s.locker.RUnlock()
	if ok {
		return errors.New("ID already exists in the connections.")
	}

	s.locker.Lock()
	s.conns[uniqueID] = conn
	s.locker.Unlock()

	return nil
}


// GetConn returns a TCP connection based on a given ID, or error if the given ID doesn't exist in the map.
func (s *Sessions) GetConn(uniqueID ConnID) (net.Conn, error) {
	s.locker.RLock()
	v, ok := s.conns[uniqueID]
	s.locker.RUnlock()
	if !ok {
		return nil, errors.New("Requested TCP connection doesn't exist.")
	}

	return v, nil
}


// DeleteConn removes a TCP connection based on a given ID.
func (s *Sessions) DeleteConn(uniqueID ConnID) {
	s.locker.Lock()
	delete(s.conns, uniqueID)
	s.locker.Unlock()
}
