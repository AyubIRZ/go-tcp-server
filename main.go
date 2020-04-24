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

	fmt.Println("TCP server started . . .")

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
			log.Println("TCP conn failed: ", err)
			continue
		}

		fmt.Println("New conn accepted! ", conn.RemoteAddr())

		wg.Add(1)
		go handleConn(conn, wg)
	}
}

func handleConn(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()

	msg, _ := bufio.NewReader(conn).ReadString('\n')
	//conn.Write([]byte("heeeeyyyyyyyyyyyy!!!!!"))

	fmt.Println(msg)
}
