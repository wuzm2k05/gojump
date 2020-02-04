package main

import (
	"bufio"
	"crypto/tls"
	//  "log"
	"net"
	//  "net/http"
	//  "fmt"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"slog"
	"time"
)

type ServerConf struct {
	JumpServerListenAddr string
}

var logger *log.Logger
var serverConf *ServerConf

func parseServerConf() {
	file, err := os.Open("serverconf.json")
	if err != nil {
		logger.Println("Warning, no serverconf.json found, will use default addr")
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
  //need ':' for tls server listen
	serverConf = &ServerConf{JumpServerListenAddr: ":4444"}
	err = decoder.Decode(&serverConf)
	if err != nil {
		logger.Println("Warning, parse serverconf.json fail, will use default addr")
	}

  //use enviroment if there is an port
  serverPort, exists := os.LookupEnv("GO_JUMP_SERVER_PORT")
  if exists {
    serverConf.JumpServerListenAddr = serverPort
  }
}

func main() {

	slog.LoggerInit("jumpServer.log")
	logger = slog.GetInstance()
	logger.Println("jumpClient start")
	parseServerConf()

	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		logger.Println("error tls")
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	ln, err := tls.Listen("tcp", serverConf.JumpServerListenAddr, config)
	if err != nil {
		logger.Println("erro listen tls")
		return
	}
	logger.Println("JumpServer listen on addr:" + serverConf.JumpServerListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {

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
		if msg[:9] == "YAEFCTqyz" {
			go handleHttpConn(conn)
			return
		} else {
			logger.Println("recieve an error message:" + msg[:9])
			conn.Write([]byte("not support!\n"))
			conn.Close()
			return
		}
	}

	//now we get the right host, need to connect with destination and inform jumpclient
	dest_conn, err := net.DialTimeout("tcp", msg[10:len(msg)-1], 10*time.Second)
	if err != nil {
		logger.Println(msg)
		logger.Println(len(msg))
		logger.Println("connect with host: " + msg[10:] + " fail!!")
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

/*
Des: this is a connection for http request.
*/
func handleHttpConn(conn net.Conn) {
	logger.Println("an HTTP connection thread")
	for {
		req, err := http.ReadRequest(bufio.NewReader(conn))
		if err != nil {
			logger.Println("recieve http request fail")
			logger.Println(err)
			return
		}
		logger.Println("recieve one http request:")
		logger.Println(req)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			logger.Println("send request to client fail")
			logger.Println(err)
			return
		}
		res.Write(conn)
		res.Body.Close()
	}
}
