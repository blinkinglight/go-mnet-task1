package main

import (
	"net"
	"testing"
	"time"
)

func init() {
	go Server(":9999")
	time.Sleep(2 * time.Second)
}

func TestMainShot0(t *testing.T) {

	nd, err := net.Dial("tcp", "127.0.0.1:9999")

	if err != nil {
		t.Fatal(err)
	}

	nd.Write([]byte("START test\n"))
	b, err := _readBytes(nd)

	if err != nil {
		t.Fatal(err)
	}

	if string(b) != "WALK zombie 1 0\n" {
		t.Fatal()
	}

	nd.Write([]byte("SHOOT 0 0\n"))

	b1, err := _readBytes(nd)

	if err != nil {
		t.Fatal(err)
	}

	if string(b1) != "BOOM test 0\n" {
		t.Fatal()
	}

	nd.Write([]byte("SHOOT 0 0\n"))

	b2, err := _readBytes(nd)

	if err != nil {
		t.Fatal(err)
	}

	if string(b2) != "BOOM test 0\n" {
		t.Fatal()
	}

	nd.Write([]byte("SHOOT 0 0\n"))

	b3, err := _readBytes(nd)

	if err != nil {
		t.Fatal(err)
	}

	if string(b3) != "BOOM test 0\n" {
		t.Fatal()
	}

	nd.Write([]byte("SHOOT 1 0\n"))

	b4, err := _readBytes(nd)

	if err != nil {
		t.Fatal(err)
	}

	if string(b4) != "BOOM test 1 zombie\n" {
		t.Fatal()
	}
}

func TestMainShot1(t *testing.T) {

	nd, err := net.Dial("tcp", "127.0.0.1:9999")

	if err != nil {
		t.Fatal(err)
	}

	nd.Write([]byte("START test\n"))
	b, err := _readBytes(nd)

	if err != nil {
		t.Fatal(err)
	}

	if string(b) != "WALK zombie 1 0\n" {
		t.Fatal()
	}

	nd.Write([]byte("SHOOT 0 0\n"))

	b1, err := _readBytes(nd)

	if err != nil {
		t.Fatal(err)
	}

	if string(b1) != "BOOM test 0\n" {
		t.Fatal()
	}
}

func TestMainShot2(t *testing.T) {

	nd, err := net.Dial("tcp", "127.0.0.1:9999")

	if err != nil {
		t.Fatal(err)
	}

	nd.Write([]byte("START test\n"))
	b, err := _readBytes(nd)

	if err != nil {
		t.Fatal(err)
	}

	if string(b) != "WALK zombie 1 0\n" {
		t.Fatal()
	}

	nd.Write([]byte("SHOOT 1 0\n"))

	b1, err := _readBytes(nd)

	if err != nil {
		t.Fatal(err)
	}

	if string(b1) != "BOOM test 1 zombie\n" {
		t.Fatalf("%s", b1)
	}
}
func _readBytes(nc net.Conn) ([]byte, error) {
	rb := make([]byte, 1024)
	n, err := nc.Read(rb)
	if err != nil {
		return nil, err
	}
	return rb[0:n], nil
}
