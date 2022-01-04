package main

import (
	"fmt"
	"syscall"
)

func main() {
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	check(err)
	fmt.Print(socket)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
