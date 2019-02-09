package main

import (
	"fmt"
	"net"
)

func main() {
	if isIPinNet("127.0.0.1", "127.20.0.0/16") {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}

func isIPinNet(ip string, network string) bool {
	_, ipnet, _ := net.ParseCIDR(network)
	
	if ipnet.Contains(net.ParseIP(ip)) {
		return true
	}
	
	return false
}
