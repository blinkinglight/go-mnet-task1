package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

var (
	flagBind = flag.String("bind", ":9999", "bind on")
)

var (
	wch = make(chan string)
)

func main() {
	flag.Parse()

	Server(*flagBind)
}

func Server(addr string) {
	nl, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := nl.Accept()
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		b := NewBoard(conn)
		b.Init()
	}
}

func NewBoard(c net.Conn) *Board {
	return &Board{
		c: c, wch: make(chan string),
		zpos: Point{0, 0},
		p:    "", alive: false,
	}
}

type Point struct {
	X int
	Y int
}

func (p *Point) Equals(pc Point) bool {
	return p.X == pc.X && p.Y == pc.Y
}

type Board struct {
	c     net.Conn
	wch   chan string
	zpos  Point
	p     string
	alive bool
}

func (b *Board) Send(s string) {
	b.wch <- s
}

func (b *Board) Init() {
	go b.Reader()
	go b.Writer()
}

func (b *Board) Reader() {
	reader := bufio.NewReader(b.c)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("%v", err)
			break
		}
		parts := strings.SplitN(s, " ", 2)
		if len(parts) < 2 {
			b.Send("Invalid command")
			continue
		}
		cmd, args := parts[0], strings.Trim(parts[1], " \n")
		switch cmd {
		case "START":
			go b.Zombie(args)
			continue
		case "SHOOT":
			var p Point
			fmt.Sscanf(args, "%d %d", &p.X, &p.Y)
			if p.Equals(b.zpos) && b.alive {
				b.alive = false
				b.Send(fmt.Sprintf("BOOM %s %d %s\n", b.p, 1, "zombie"))
			} else {
				b.Send(fmt.Sprintf("BOOM %s 0\n", b.p))
			}
		}
	}
}

func (b *Board) Writer() {
	for {
		select {
		case l, ok := <-b.wch:
			if !ok {
				return
			}
			_, err := b.c.Write([]byte(l))
			if err != nil {
				log.Printf("%v", err)
				return
			}
		}
	}
}

func (b *Board) Zombie(p string) {
	b.p = p
	b.alive = true
	for b.alive {
		b.zpos.X++
		b.Send(fmt.Sprintf("WALK zombie %d %d\n", b.zpos.X, b.zpos.Y))
		time.Sleep(2 * time.Second)
	}
}
