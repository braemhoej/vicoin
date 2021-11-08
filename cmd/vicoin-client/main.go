package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"vicoin/internal/registration"
)

func getExternalIP() string {
	apiUrl := "https://api.ipify.org?format=text"
	response, err := http.Get(apiUrl)
	if err != nil {
		log.Panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panic(err)
		}
	}(response.Body)
	externalIp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Panic(err)
	}
	return string(externalIp)
}

func getString() string {
	msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	} else {
		return strings.TrimSpace(msg)
	}
}
func getAddressFromUser() net.TCPAddr {
	fmt.Println("Please input target IP-address: (none to initialise new network)")
	ip := getString()
	var port = ""
	if ip != "" {
		fmt.Println("Please input port-number:")
		port = getString()
	}
	portInt, _ := strconv.Atoi(port)
	target := net.TCPAddr{
		IP:   net.ParseIP(ip),
		Port: portInt,
	}
	return target
}

func main() {
	registration.RegisterStructsWithGob()
	fmt.Println("External IP: " + getExternalIP())
}
