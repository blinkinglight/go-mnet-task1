package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync/atomic"
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
		NewBoard(conn)
	}
}

func NewBoard(c net.Conn) *Board {
	b := &Board{
		c: c, wch: make(chan string),
		zpos: Point{0, 0},
		p:    "", alive: false,
		points:  make(map[string]int),
		running: 0,
	}
	b.initContext()
	go b.Reader()
	return b
}

type Point struct {
	X int
	Y int
}

func (p *Point) Equals(pc Point) bool {
	return p.X == pc.X && p.Y == pc.Y
}

type Board struct {
	c       net.Conn
	wch     chan string
	zpos    Point
	p       string
	alive   bool
	points  map[string]int
	running int32
	ctx     context.Context
	cfn     context.CancelFunc
}

func (b *Board) initContext() {
	b.ctx, b.cfn = context.WithCancel(context.Background())
}

func (b *Board) Send(s string) error {
	_, err := b.c.Write([]byte(s))
	return err
}

func (b *Board) Reader() {
	reader := bufio.NewReader(b.c)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			if b.cfn != nil {
				b.cfn()
			}
			break
		}

		parts := strings.SplitN(s, " ", 2)
		if len(parts) < 2 {
			log.Printf("command to short")
			continue
		}

		// rename and strip things
		cmd, args := parts[0], strings.Trim(parts[1], " \n")

		switch cmd {
		case "START":
			go b.Zombie(args)
		case "SHOOT":
			b.Shoot(args)
		}
	}
}

func (b *Board) Shoot(args string) {
	if atomic.LoadInt32(&b.running) == 0 {
		return
	}

	var p Point
	fmt.Sscanf(args, "%d %d", &p.X, &p.Y)
	if p.Equals(b.zpos) && b.alive {
		b.alive = false
		b.cBoom(true)
		return
	}
	b.cBoom(false)
}

func (b *Board) Zombie(p string) {
	if !atomic.CompareAndSwapInt32(&b.running, 0, 1) {
		return
	}
	b.initContext()
	b.p = p
	if _, ok := b.points[p]; !ok {
		b.points[p] = 0
	}
	b.alive = true
	b.zpos = Point{0, 0}
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ticker.C:
			b.zpos.X++
			b.cWalkTo(b.zpos)
		case <-b.ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (b *Board) cBoom(s bool) error {
	if s && atomic.CompareAndSwapInt32(&b.running, 1, 0) {
		b.cfn()
		b.points[b.p]++
		return b.Send(fmt.Sprintf("BOOM %s %d %s\n", b.p, b.points[b.p], "zombie"))
	}
	return b.Send(fmt.Sprintf("BOOM %s %d\n", b.p, b.points[b.p]))
}

func (b *Board) cWalkTo(p Point) {
	b.Send(fmt.Sprintf("WALK zombie %d %d\n", p.X, p.Y))
}
