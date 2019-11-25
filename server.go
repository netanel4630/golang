package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
) 

const (
    SEND_TO_SERVER = "1"
    SEND_TO_CLIENT = "2"
)

type client struct {
    address string
    c net.Conn
    isAlive bool
}

const INTERVAL_PRINTS = 5
const PROTOCOL = "tcp"
const PORT = ":8081"
var clients []client //list for clients connection info

/************************************************************************
* Function: addClientToList()
* Purpose:  Add client connection info to list
* Input:    c - generic stream-oriented network connection  
* Return:   None
************************************************************************/
func initClient(c net.Conn) client{
	temp := client{address: c.RemoteAddr().String()}
	temp.c = c
	temp.isAlive = true

	return temp
}

/************************************************************************
* Function: addClientToList()
* Purpose:  Add client connection info to list
* Input:    c - generic stream-oriented network connection  
* Return:   None
************************************************************************/
func addClientToList(c client) {
	clients = append(clients, c)
}

/************************************************************************
* Function: editStillAlive()
* Purpose:  change isAlive status
* Input:    address - Client address
* Return:   None
************************************************************************/
func editStillAlive(address string, aliveStatus bool) {
	for i := 0; i < len(clients); i++{
		if clients[i].c.RemoteAddr().String() == address {
			clients[i].isAlive = aliveStatus
			break
		}
	}
}

/************************************************************************
* Function: handleConnection()
* Purpose:  To handle messages from clients and add thier connection info
*			to list
* Input:    c - generic stream-oriented network connection  
* Return:   None
************************************************************************/
func handleConnection(c net.Conn) {

	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	addClientToList(initClient(c))
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			editStillAlive(c.RemoteAddr().String(), false)
			return
		}		
		temp := strings.TrimSpace(string(netData))
		//terminate message from client
		if temp[1:] == "STOP" {
			fmt.Println(temp[1:])
			break
		}
		if msg_parser(temp) == true {
			length := len(temp)
			result := "ACK " + strconv.Itoa(length - 1) + "\n"
			c.Write([]byte(string(result)))	
		} else{
			result := "NACK " + "\n"
			c.Write([]byte(string(result)))	
		}		
	}
	editStillAlive(c.RemoteAddr().String(), false)
	c.Close()
}

/************************************************************************
* Function: printClientsStatus()
* Purpose:  Print the status and clients
* Input:    addresses - adresses of clients and thier status 
* Return:   None
************************************************************************/
func printClientsStatus(addresses []string, status string){
	fmt.Println("\n\n")	
	fmt.Println(status)
	for i := 0; i < len(addresses); i++ {
		fmt.Println(addresses[i])
	}
}

/************************************************************************
* Function: checkConnection()
* Purpose:  To check wich clients still connected and print it 
* Input:    None 
* Return:   None
************************************************************************/
func checkConnection(){
	var conn []string
	var notConn []string
	var temp []string

	for {
		for i := 0; i < len(clients); i++ {
			if clients[i].isAlive == true {
				conn = append(conn, clients[i].c.RemoteAddr().String())			
			} else {
				notConn = append(notConn, clients[i].c.RemoteAddr().String())	
			}
		}
		printClientsStatus(conn, "Active Clients")
		printClientsStatus(notConn, "Non-Active Clients")
		conn = temp
		notConn = temp
		time.Sleep(time.Second * INTERVAL_PRINTS)	
	}
}

/************************************************************************
* Function: msg_parser()
* Purpose:  Parser Message
* Input:    msg - message from client 
* Return:   None
************************************************************************/
func msg_parser(msg string) bool{
	var index int
	var clientAdd string
	var clientMSG string

	if msg[:1] == SEND_TO_CLIENT{
		index = strings.Index(msg, ";")
		clientAdd = msg[1 : index]
		clientMSG = msg[index + 1:]
		return sendMsgToClient(clientAdd, clientMSG)
 	} else if msg[:1] == SEND_TO_SERVER{
 		fmt.Println(msg[1:])
 		return true
 	} else {
 		return false
 	}
}

/************************************************************************
* Function: sendMsgToClient()
* Purpose:  Send message from client to another
* Input:    clientAdd - client adress to send
*			clientMSG = message to send
* Return:   None
************************************************************************/
func sendMsgToClient(clientAdd string, clientMSG string) bool {
	for i := 0; i < len(clients); i++ {
		if clients[i].c.RemoteAddr().String() == clientAdd{
			sendClient := "ACK " + clientMSG + "\n"
			clients[i].c.Write([]byte(string(sendClient)))
			return true
		}
	}
	return false
}

func main() {
	l, err := net.Listen(PROTOCOL, PORT)
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