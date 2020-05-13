package session

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"github.com/google/uuid"
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
func (s *Sessions) AddConn(conn net.Conn) (ConnID, error) {
	s.locker.RLock()
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	uniqueID := ConnID(id.String())
	_, ok := s.conns[uniqueID]
	s.locker.RUnlock()
	if ok {
		return "", errors.New("ID already exists in the connections.")
	}

	s.locker.Lock()
	s.conns[uniqueID] = conn
	s.locker.Unlock()

	return uniqueID, nil
}


// GetConn returns a TCP connection based on a given ID, or error if the given ID doesn't exist in the map.
func (s *Sessions) GetConn(uniqueID ConnID) (net.Conn, error) {
	s.locker.RLock()
	conn, ok := s.conns[uniqueID]
	s.locker.RUnlock()
	if !ok {
		return nil, errors.New("Requested TCP connection doesn't exist.")
	}

	return conn, nil
}


// DeleteConn removes a TCP connection based on a given ID.
func (s *Sessions) DeleteConn(uniqueID ConnID) {
	s.locker.Lock()
	delete(s.conns, uniqueID)
	s.locker.Unlock()
}
