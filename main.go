package main

import (
	"bufio"
	"fmt"
	"github.com/ayubirz/go-tcp-server/pkg/session"
	"log"
	"net"
	"sync"
)

// general constants are here
const (
	TCPPort       = ":6060"
	authDelimiter = "AUTH="
	infoConnSuccessful = "INFO_CONN_SUCCESSFUL"
)

// All TCP conn related errors are defined here.
const (
	errorIDExists 		= "ERROR_ID_EXISTS"
	errorIDNotSpecified = "ERROR_ID_NOT_SPECIFIED"
)

func main() {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	sessions := session.Init()

	fmt.Println("* * * TCP server started * * *")

	wg.Add(1)
	go handleListener(wg, sessions)
}

func handleListener(wg *sync.WaitGroup, sessions *session.Sessions) {
	defer wg.Done()

	listener, err := net.Listen("tcp", TCPPort)
	if err != nil {
		log.Fatal("TCP listener error: ", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("TCP connection failed: ", err)
			continue
		}

		fmt.Println("New connection accepted! ", conn.RemoteAddr())

		wg.Add(1)
		go handleConn(conn, wg, sessions)
	}
}

func handleConn(conn net.Conn, wg *sync.WaitGroup, sessions *session.Sessions) {
	defer wg.Done()
	defer func() {
		_ = conn.Close()
	}()

	buf := bufio.NewReader(conn)

	msg, err := buf.ReadString('\n')
	if err != nil {
		log.Fatal("TCP connection read failed: ", err)
	}

	if msg[:len(authDelimiter)] != authDelimiter{
		_, _ = fmt.Fprintln(conn, errorIDNotSpecified)
		return
	}

	connID := session.ConnID(msg[len(authDelimiter) : len(msg)-1])
	if err := sessions.AddConn(connID, conn); err != nil {
		_, _ = fmt.Fprintln(conn, errorIDExists)
		return
	}
	_, _ = fmt.Fprintln(conn, infoConnSuccessful)

	for {
		msg, err = buf.ReadString('\n')
		if err != nil {
			log.Println("TCP connection read failed: ", err)
			break
		}
		fmt.Println("<client>: ", msg)
		_, _ = fmt.Fprintln(conn, "Hi from TCP server :)")
	}
}
