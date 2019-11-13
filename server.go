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

var l []net.Conn //list for clients connection info

/************************************************************************
* Function: handleConnection()
* Purpose:  To handle messages from clients and add thier connection info
*			to list
* Input:    c - generic stream-oriented network connection  
* Return:   None
************************************************************************/
func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		l = append(l, c)

		temp := strings.TrimSpace(string(netData))
		//terminate message from client
		if temp == "STOP" {
			break
		}
		//still alive message from client
		else if temp == "OK" { 

		}
		//any other message from client
		else {
			length := len(temp)
			result := "ACK " + strconv.Itoa(length) + "\n"
			c.Write([]byte(string(result)))		
		}

	}
	c.Close()
}

/************************************************************************
* Function: checkConnection()
* Purpose:  To check wich clients still connected
* Input:    None 
* Return:   None
************************************************************************/
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