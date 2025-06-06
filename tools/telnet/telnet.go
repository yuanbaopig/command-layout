package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	conn        net.Conn
	isConnected bool
	mu          sync.Mutex
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		mu.Lock()
		connected := isConnected
		mu.Unlock()

		if !connected {
			// 命令模式
			fmt.Print("telnet> ")
			if !scanner.Scan() {
				break
			}
			handleCommand(scanner.Text())
		} else {
			// 交互模式
			fmt.Print("> ")
			if !scanner.Scan() {
				break
			}
			handleInteractiveInput(scanner.Text())
		}
	}
}

func handleCommand(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "open":
		if len(parts) < 3 {
			fmt.Println("Usage: open <host> <port>")
			return
		}
		host := parts[1]
		port := parts[2]
		connect(host, port)
	case "quit":
		fmt.Println("Exiting...")
		os.Exit(0)
	default:
		fmt.Println("Unknown command. Supported: open, quit")
	}
}

func connect(host, port string) {
	c, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		fmt.Printf("Error connecting to %s:%s - %v\n", host, port, err)
		return
	}

	mu.Lock()
	conn = c
	isConnected = true
	mu.Unlock()

	fmt.Printf("Connected to %s:%s\n", host, port)
	go readServerResponses()
}

func readServerResponses() {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	mu.Lock()
	defer mu.Unlock()
	if isConnected {
		fmt.Println("\nConnection closed by server")
		conn.Close()
		isConnected = false
	}
}

func handleInteractiveInput(input string) {
	mu.Lock()
	defer mu.Unlock()

	if strings.TrimSpace(input) == "close" {
		fmt.Println("Closing connection...")
		conn.Close()
		isConnected = false
		return
	}

	if _, err := fmt.Fprintln(conn, input); err != nil {
		fmt.Printf("Error sending data: %v\n", err)
		conn.Close()
		isConnected = false
	}
}
