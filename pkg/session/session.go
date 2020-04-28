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


// AddConn adds a new connection to the sessions map.
func (b *Sessions) AddConn(uniqueID ConnID, conn net.Conn) error {
	b.locker.RLock()
	_, ok := b.conns[uniqueID]
	b.locker.RUnlock()
	if ok {
		return errors.New("ID already exists in the connections.")
	}

	b.locker.Lock()
	b.conns[uniqueID] = conn
	b.locker.Unlock()

	return nil
}


// GetConn returns a TCP connection based on a given ID, or error if the given ID doesn't exist in the map.
func (b *Sessions) GetConn(uniqueID ConnID) (net.Conn, error) {
	b.locker.RLock()
	v, ok := b.conns[uniqueID]
	b.locker.RUnlock()
	if ok && v != nil {
		return v, nil
	}

	return nil, errors.New("Requested TCP connection doesn't exist.")
}


// DeleteConn removes a TCP connection based on a given ID.
func (b *Sessions) DeleteConn(uniqueID ConnID) {
	b.locker.Lock()
	delete(b.conns, uniqueID)
	b.locker.Unlock()
}
