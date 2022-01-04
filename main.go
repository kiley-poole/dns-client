package main

import (
	"syscall"
)

func main() {
	sa := &syscall.SockaddrInet4{}
	googleDns := &syscall.SockaddrInet4{
		Addr: [4]byte{8, 8, 8, 8},
		Port: 53,
	}
	var data []byte
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	check(err)
	err = syscall.Bind(socket, sa)
	check(err)
	err = syscall.Sendto(socket, data, syscall.MSG_CONFIRM, googleDns)
	check(err)
	syscall.Close(socket)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
