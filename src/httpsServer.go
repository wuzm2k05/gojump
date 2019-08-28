package main

import (
  "bufio"
  "crypto/tls"
//  "log"
  "net"
//  "net/http"
  "fmt"
  "time"
  "io"
)

func main(){
  cert,err := tls.LoadX509KeyPair("server.crt", "server.key")
  if err != nil {
    fmt.Println("error tls")
    return
  }

  config := &tls.Config{Certificates: []tls.Certificate{cert}}
  ln, err := tls.Listen("tcp", ":4444", config)
  if err != nil {
    fmt.Println("erro listen tls")
    return
  }

  for {
    conn, err := ln.Accept()
    if err != nil {
      fmt.Println(err)
      continue
    }
    go handleConn(conn)
  }
}

func handleConn(conn net.Conn){
  //defer conn.Close()
  //expect magic connection here
  r := bufio.NewReader(conn)
  msg, err := r.ReadString('\n')
  if err != nil {
    fmt.Println("read error")
    fmt.Println(err)
    return
  }
  
  //check if the message start with magic string
  if msg[:9] != "XAEFCTqyz" {
    fmt.Println("recieve an error message:"+msg[:9])
    conn.Write([]byte("not support!\n"))
    conn.Close()
    return
  }

  //now we get the right host, need to connect with destination and inform jumpclient
  dest_conn, err := net.DialTimeout("tcp", msg[9:],10*time.Second)
  if err != nil {
    conn.Write([]byte("not support!\n"))
    return
  }
  conn.Write([]byte("okay"))

  go transfer(dest_conn, conn)
  go transfer(conn, dest_conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
    defer destination.Close()
    defer source.Close()
    io.Copy(destination, source)
}
