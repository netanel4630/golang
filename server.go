package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var l []net.Conn //list for clients info
////
func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println(netData)

		l = append(l, c)

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}
		length := len(temp)
		result := "ACK " + strconv.Itoa(length) + "\n"
		c.Write([]byte(string(result)))
	}
	c.Close()
}


func checkConnection(){

	for {
		for i := 0; i < len(l); i++{
			fmt.Println(l[i].RemoteAddr().String())
		}
		time.Sleep(time.Second * 5)	
	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	//call the goroutine  to check connection
	go checkConnection();

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}