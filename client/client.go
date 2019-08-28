package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {
	nc, err := net.Dial("tcp", ":9999")
	if err != nil {
		panic(err)
	}
	go func() {
		reader := bufio.NewReader(nc)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			fmt.Print(s)
		}
	}()

	println("START c1")
	nc.Write([]byte("START c1\n"))
	time.Sleep(500 * time.Millisecond)
	println("SHOOT 0 0")
	nc.Write([]byte("SHOOT 0 0\n"))
	println("SHOOT 0 0")
	nc.Write([]byte("SHOOT 0 0\n"))
	println("SHOOT 0 0")
	nc.Write([]byte("SHOOT 0 0\n"))
	time.Sleep(2000 * time.Millisecond)
	println("SHOOT 2 0")
	nc.Write([]byte("SHOOT 2 0\n"))

	time.Sleep(1 * time.Second)
}
