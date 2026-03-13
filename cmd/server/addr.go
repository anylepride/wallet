package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

func getGrpcAddr() string {
	var (
		port   int
		ipAddr string
		err    error
	)

	ipAddr, err = getLocalAddr()
	if err != nil {
		log.Fatalf("getLocalAddr err, %v", err)
	}
	port, err = getAvailablePort()
	if err != nil {
		log.Fatalf("getAvailablePort err, %v", err)
	}
	return net.JoinHostPort(ipAddr, strconv.Itoa(port))
}

func getLocalAddr() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no ipv4 address found")
}

func getAvailablePort() (int, error) {
	address, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if err != nil {
		return 0, err
	}
	ln, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}
	defer func(ln *net.TCPListener) {
		_ = ln.Close()
	}(ln)
	return ln.Addr().(*net.TCPAddr).Port, nil
}
