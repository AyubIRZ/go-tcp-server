package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

const TCPAddr = "localhost:6060"

func main() {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	fmt.Println("* * * TCP server started * * *")

	wg.Add(1)
	go handleListener(wg)
}

func handleListener(wg *sync.WaitGroup) {
	defer wg.Done()

	listener, err := net.Listen("tcp", TCPAddr)
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
		go handleConn(conn, wg)
	}
}

func handleConn(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	buf := bufio.NewReader(conn)
	for {
		msg, err := buf.ReadString('\n')
		if err != nil {
			log.Println("TCP connection read failed: ", err)
			break
		}
		fmt.Println("<client>: ", msg)
		_, _ = fmt.Fprint(conn, "Hi from TCP server :)\n")
	}
	_ = conn.Close()
}
