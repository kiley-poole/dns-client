package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"syscall"
)

func main() {
	socketAddr := &syscall.SockaddrInet4{}
	googleDns := &syscall.SockaddrInet4{
		Addr: [4]byte{8, 8, 8, 8},
		Port: 53,
	}
	DNSMessage := buildDNSMessage("google.com")

	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	check(err)

	err = syscall.Bind(socket, socketAddr)
	check(err)

	err = syscall.Sendto(socket, DNSMessage, 0, googleDns)
	check(err)

	res := make([]byte, 512)

	_, _, err = syscall.Recvfrom(socket, res, 0)
	check(err)

	parseAndPrint(res)

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

type DNSHeader struct {
	ID      uint16
	Opt1    byte
	Opt2    byte
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

func parseAndPrint(res []byte) {
	var resHeader DNSHeader
	buf := bytes.NewReader(res)

	err := binary.Read(buf, binary.BigEndian, &resHeader)
	check(err)

	host := decodeQuestionData(buf)
	ip := decodeAnswerData(buf)
	fmt.Printf("HOST NAME: %s IP ADDRESS: %d.%d.%d.%d", host, ip[0], ip[1], ip[2], ip[3])

}

type AnswerHeader struct {
	Offset uint16
	Typ    uint16
	Class  uint16
	TTL    uint32
	Rdlen  uint16
}

func decodeAnswerData(buf *bytes.Reader) []byte {
	var ansHead AnswerHeader
	err := binary.Read(buf, binary.BigEndian, &ansHead)
	check(err)

	ip := make([]byte, ansHead.Rdlen)
	err = binary.Read(buf, binary.BigEndian, &ip)
	check(err)

	return ip
}

func decodeQuestionData(buf *bytes.Reader) string {
	var (
		err    error
		qType  uint16
		qClass uint16
		host   []string
	)
	for {
		var size byte
		err = binary.Read(buf, binary.BigEndian, &size)
		check(err)
		if size == 0 {
			break
		}
		e := make([]byte, size)
		err = binary.Read(buf, binary.BigEndian, &e)
		check(err)
		host = append(host, string(e))
	}

	err = binary.Read(buf, binary.BigEndian, &qType)
	check(err)

	err = binary.Read(buf, binary.BigEndian, &qClass)
	check(err)

	return strings.Join(host, ".")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
