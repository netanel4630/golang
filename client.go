package main

import "net"
import "fmt"
import "bufio"
import "os"
import "strings"
import "time"

const (
    SEND_TO_SERVER = "1"
    SEND_TO_CLIENT = "2"
    INTERVAL_PRINTS = 100
)

/************************************************************************
* Function: receiveMsg()
* Purpose:  Receive message from Server
* Input:    conn - connection info
* Return:   None
************************************************************************/
func receiveMsg(conn net.Conn) {
  for{
    message, _ := bufio.NewReader(conn).ReadString('\n')
    if message[:4] == "NACK"{
      fmt.Print("Client not exist or disconnected")
    } else {
      fmt.Print("\n")
      fmt.Print(message[4:])
      fmt.Print("1 for send to Server, 2 for send to Client: ")
      //fmt.Print("Message from server: "+message)  -> I think its better format, but I follow your order
    }

  }
}


func main() {
  var temp string
  // connect to this socket
  conn, _ := net.Dial("tcp", "127.0.0.1:8081")
  go receiveMsg(conn)
  fmt.Print("1 for send to Server, 2 for send to Client: ")
  for {
    // read in input from stdin
    reader := bufio.NewReader(os.Stdin)
    text, _ := reader.ReadString('\n')
    if text == SEND_TO_SERVER + "\n"{
      reader := bufio.NewReader(os.Stdin)
      fmt.Print("Text to send: ")
      text, _ := reader.ReadString('\n')
      // send to socket
      fmt.Fprintf(conn, SEND_TO_SERVER + text + "\n")
    } else if text == SEND_TO_CLIENT + "\n" {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Client address: ")
        text, _ := reader.ReadString('\n')
        temp = text
        reader = bufio.NewReader(os.Stdin)
        temp += ";"
        fmt.Print("Text to send:")
        text, _ = reader.ReadString('\n')
        temp += text
        // send to socket
        temp = strings.ReplaceAll(temp, "\n","")
        fmt.Fprintf(conn, SEND_TO_CLIENT + temp + "\n")
        temp = strings.ReplaceAll(temp, ";","")
        fmt.Println(temp[1:] + "\n")
    } else {
        fmt.Println("Error! please enter 1 or 2")
        fmt.Print("1 for send to Server, 2 for send to Client: ")
        continue
    }   

    time.Sleep(time.Millisecond * INTERVAL_PRINTS)
  }
}