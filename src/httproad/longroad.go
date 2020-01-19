package httproad

import (
	"bufio"
	"crypto/tls"
	"net"
	"net/http"
	//"os"
	"log"
	"slog"
	"sync"
)

var reqchan = make(chan *http.Request)
var reschan = make(chan *http.Response)
var once sync.Once
var server_conn net.Conn
var logger *log.Logger

/*
Client must call resp.Body.Close when finished reading resp.Body
This function should be thread safe
This function send the msg through channel req/res whose length is 1,
to serialize all req/res.
*/
func SendHttpReq(url string, req *http.Request) *http.Response {
	once.Do(func() {
		logger = slog.GetInstance()
		go httpMsgRoadThread(url)
	})
	reqchan <- req
	return <-reschan
}

/*
this thread run forever, read request from reqchan,
then send the request to remote server,
get response from remote server,
return response to reschan.
Why we need this thread? to serialize all http req
*/
func httpMsgRoadThread(url string) {
	for {
		req := <-reqchan
		conn := getTlsConn(url)
		req.Write(conn)
		res, ok := http.ReadResponse(bufio.NewReader(conn), req)
		if ok != nil {
			logger.Println("error get Response")
		} else {
			logger.Println("get response correct")
		}
		reschan <- res
	}
}

/*
this should be an indpendent pacakge later.
*/
func getTlsConn(url string) net.Conn {
	if server_conn != nil {
		return server_conn
	}

	var err error

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	server_conn, err = tls.Dial("tcp", url, conf)
	if err != nil {
		logger.Println("error estabish connection with server\n")
		logger.Println(err)
		return nil
	}

	// so send the magic number to server, then server can identify
	// this is a connection for http not https
	_, err = server_conn.Write([]byte("YAEFCTqyz")) // magic number for http connection
	if err != nil {
		logger.Println("jumpClient: error write")
		logger.Println(err)
		server_conn.Close()
		return nil
	}

	//do we need response from server, may not now

	return server_conn
}
