package main

import (
  "bufio"
  "crypto/tls"
//  "log"
  "net"
//  "net/http"
//  "fmt"
  "time"
  "io"
  "slog"
  "log"
)

var logger *log.Logger

func main(){
  
  slog.LoggerInit("jumpClient.log")
  logger = slog.GetInstance()
  logger.Println("jumpClient start")

  cert,err := tls.LoadX509KeyPair("server.crt", "server.key")
  if err != nil {
    logger.Println("error tls")
    return
  }

  config := &tls.Config{Certificates: []tls.Certificate{cert}}
  ln, err := tls.Listen("tcp", ":4444", config)
  if err != nil {
    logger.Println("erro listen tls")
    return
  }

  for {
    conn, err := ln.Accept()
    if err != nil {
      logger.Println(err)
      continue
    }
    go handleConn(conn)
  }
}

func handleConn(conn net.Conn){

  logger.Println("one connection request recieved\n")
  //defer conn.Close()
  //expect magic connection here
  r := bufio.NewReader(conn)
  msg, err := r.ReadString('\n')
  if err != nil {
    logger.Println("read error")
    logger.Println(err)
    conn.Close()
    return
  }
  
  //check if the message start with magic string
  if msg[:9] != "XAEFCTqyz" {
    logger.Println("recieve an error message:"+msg[:9])
    conn.Write([]byte("not support!\n"))
    conn.Close()
    return
  }

  //now we get the right host, need to connect with destination and inform jumpclient
  dest_conn, err := net.DialTimeout("tcp", msg[10:len(msg)-1],10*time.Second)
  if err != nil {
    logger.Println(msg)
    logger.Println(len(msg))
    logger.Println("connect with host: " + msg[10:] +" fail!!")
    conn.Write([]byte("not support!\n"))
    conn.Close()
    logger.Println(err)
    return
  }
  //send response to jumpClient
  conn.Write([]byte("okay"))

  go transfer(dest_conn, conn)
  go transfer(conn, dest_conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
    defer destination.Close()
    defer source.Close()
    io.Copy(destination, source)
}
