package main

import (
	"strings"
	"syscall"
)

func main() {
	sa := &syscall.SockaddrInet4{}
	googleDns := &syscall.SockaddrInet4{
		Addr: [4]byte{8, 8, 8, 8},
		Port: 53,
	}

	DNSMessage := buildDNSMessage("google.com")

	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	check(err)

	err = syscall.Bind(socket, sa)
	check(err)

	err = syscall.Sendto(socket, DNSMessage, syscall.MSG_DONTWAIT, googleDns)
	check(err)

	err = syscall.Close(socket)
	check(err)
}

func buildDNSMessage(host string) []byte {
	dnsHeader := []byte{0xFF, 0xFF, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	question := buildQuestion(host)

	message := append(dnsHeader, question...)

	return message
}

func buildQuestion(host string) []byte {
	var question []byte
	elements := strings.Split(host, ".")

	for _, e := range elements {
		l := len(e)
		question = append(question, byte(l))
		question = append(question, e...)
	}
	question = append(question, []byte{0x00, 0x00, 0x01, 0x00, 0x01}...)

	return question
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
