package main

import (
	"bufio"
	"fmt"
	"github.com/ayubirz/go-tcp-server/pkg/session"
	"log"
	"net"
	"net/url"
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

	connID, err := sessions.AddConn(conn)
	if err != nil {
		_, _ = fmt.Fprintln(conn, errorIDExists)
		return
	}
	defer sessions.DeleteConn(connID)
	_, _ = fmt.Fprintln(conn, infoConnSuccessful)

	for {
		msg, err := buf.ReadString('\n')
		if err != nil {
			log.Println("TCP connection read failed: ", err)
			return
		}

		msg, _ = url.QueryUnescape(msg)

		fmt.Println("<client>: ", msg)
		_, err = fmt.Fprintln(conn, "Hi from TCP server :)")
		if err != nil {
			fmt.Println("some write error: ", err.Error())
		}
	}
}
