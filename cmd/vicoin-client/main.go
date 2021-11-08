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
	"vicoin/crypto"
	"vicoin/internal/account"
	"vicoin/internal/client"
	"vicoin/internal/node"
	"vicoin/internal/registration"
	"vicoin/network"
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
func getAddressFromUser() *net.TCPAddr {
	fmt.Println("Please input target IP-address: ")
	fmt.Print(" >   ")
	ip := net.ParseIP(getString())
	if ip == nil {
		fmt.Println("Cancelling : Invalid IP entered")
		return nil
	}
	fmt.Println("Please input port-number:")
	fmt.Print(" >   ")
	port, err := strconv.Atoi(getString())
	if err != nil || port > 65535 || port < 1 {
		fmt.Println("Cancelling : Invalid port entered")
		return nil
	}
	target := &net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	return target
}

func getCredentialsFromUser() (public *crypto.PublicKey, private *crypto.PrivateKey) {
	for {
		fmt.Println("Please provide public key: ")
		var err error
		fmt.Print(" >   ")
		public, err = new(crypto.PublicKey).FromString(getString())
		if err != nil {
			fmt.Println("Error when parsing public key : ", err)
			continue
		}
		break
	}
	for {
		fmt.Println("Please provide private key: ")
		var err error
		fmt.Print(" >   ")
		private, err = new(crypto.PrivateKey).FromString(getString())
		if err != nil {
			fmt.Println("Error when parsing private key : ", err)
			continue
		}
		break
	}
	return public, private
}

func createAndConfigureClient() (*client.Client, error) {
	fmt.Println("Configuring client ...")
	socketToNode := make(chan interface{})
	nodeToClient := make(chan account.SignedTransaction)
	dialer, err := network.NewTCPDialer()
	if err != nil {
		return nil, err
	}
	listener, err := network.NewTCPListener()
	if err != nil {
		return nil, err
	}
	socket := network.NewPolysocket(socketToNode, dialer, listener)
	node, err := node.NewNode(socket, socketToNode, nodeToClient)
	if err != nil {
		return nil, err
	}
	ledger := account.NewLedger()
	client, err := client.NewClient(ledger, node, nodeToClient)
	if err != nil {
		return nil, err
	}
	fmt.Println("Client configured ...")
	return client, nil
}

func login(client *client.Client) {
	for {
		fmt.Println("Provide credentials (P), or generate new (G)? :")
		fmt.Print(" >   ")
		answer := getString()
		switch strings.ToUpper(answer) {
		case "P":
			client.ProvideCredentials(getCredentialsFromUser())
			return
		case "G":
			public, private, err := crypto.KeyGen(2048)
			if err != nil {
				fmt.Println("Error generating credentials :  ", err)
				continue
			}
			fmt.Println("Credentials successfully generated")
			client.ProvideCredentials(public, private)
			return
		default:
			fmt.Println("Invalid input")
			continue
		}
	}
}

func printHelp() {
	fmt.Println("Valid commands: ")
	format := "%-12s : %-12s\n"
	//TODO: Implement
	fmt.Printf(format, "quit", "exits the shell")
	fmt.Printf(format, "help", "prints this")
	fmt.Printf(format, "connect", "initiate shell interaction for establishing TCP connection")
	fmt.Printf(format, "transfer", "initiate shell interaction for transfering funds")
	fmt.Printf(format, "balance", "initiate shell interaction for looking up balance")
}

func handleInput(input string, client *client.Client) (quit bool) {
	switch input {
	case "help":
		printHelp()
	case "connect":
		address := getAddressFromUser()
		err := client.Connect(address)
		if err != nil {
			fmt.Println("Error : ", err)
		} else {
			fmt.Println("Successfully connected to : " + address.String())
		}
	case "balance":
		fmt.Println("Please enter account to look up balance: (Nothing for own account)")
		account := getString()
		balance := ""
		switch account {
		case "":
			balance = fmt.Sprintf("%f", client.GetBalance(client.GetAccount()))
		default:
			balance = fmt.Sprintf("%f", client.GetBalance(account))
		}
		fmt.Println("Balance : " + balance)
	case "transfer":
		fmt.Println("Please enter recipient account (Nothing to cancel)")
		account := getString()
		switch account {
		case "":
			fmt.Println("Transaction cancelled")
			return false
		default:
			fmt.Println("Transferring to # " + account)
			fmt.Println("Please enter amount (Nothing to cancel)")
			amount := getString()
			switch amount {
			case "":
				fmt.Println("Transaction cancelled")
				return false
			default:
				if floatAmount, err := strconv.ParseFloat(amount, 64); err == nil {
					err := client.Transfer(floatAmount, account)
					if err != nil {
						fmt.Println("Error performing transaction : ", err)
					} else {
						fmt.Println("Transaction successfully performed")
					}
				} else {
					fmt.Println("Error NAN")
				}
			}
		}
	case "quit":
		errs := client.Close()
		fmt.Println("Quitting ... ")
		for err := range errs {
			fmt.Println("Error closing connection : ", err)
		}
		return true
	default:
		fmt.Println("Illegal input: (Enter 'help' for list of commands)")
	}
	return false
}

func main() {
	registration.RegisterStructsWithGob()
	client, err := createAndConfigureClient()
	if err != nil {
		fmt.Println("Fatal error: ", err)
		return
	}
	fmt.Println("Listening at IP: " + getExternalIP() + " : " + client.GetPort())
	login(client)
	fmt.Println("\nEnter 'help' for list of commands")
	for {
		fmt.Print(" >   ")
		input := getString()
		quit := handleInput(input, client)
		if quit {
			break
		}
	}
}
