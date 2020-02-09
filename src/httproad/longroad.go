package httproad

import (
	"bufio"
	//"crypto/tls"
	"net"
	"net/http"
	//"os"
	"log"
	"slog"
	"sync"
	"time"
)

var reqchan = make(chan *http.Request)
var reschan = make(chan *http.Response)
var once sync.Once
var server_conn net.Conn
var logger *log.Logger
var startTimeoutChan = make(chan bool)
var timeoutChan = make(chan uint)
var bodyDoneChan = make(chan uint, 4) //2 should be okay, but we use 4 anyway
var msgId uint  //intitliazed to zero


/*
Client must call resp.Body.Close when finished reading resp.Body
This function should be thread safe
This function send the msg through channel req/res whose length is 1,
to serialize all req/res.
*/
func SendHttpReq(url string, req *http.Request) (*http.Response,uint) {
	once.Do(func() {
		logger = slog.GetInstance()
		go httpMsgRoadThread(url)
		go timeoutThread()
	})
	reqchan <- req
	res := <- reschan
	logger.Println("get response from thread")
	return res, msgId
}

/*
client must call RecResBodyDone after it think the content in the body is transfered complete. 
Otherwise httpRoad will close the connection with jumpServer to make sure client readBody will return
after timeout.
*/
func RecResBodyDone(msgId uint){
 bodyDoneChan <- msgId
}

/*
this go routine run forever, to close connection if timeout
*/
func timeoutThread(){
  for {
    <- startTimeoutChan
    go func(id uint){
      time.Sleep(10 * time.Second)
      timeoutChan <- id
    }(msgId)
    for {
	 
	    select {
	      case expireId := <-timeoutChan:
		if expireId == msgId {
		  //the timer expired, so need to close the connection
		  server_conn.Close()
		  server_conn = nil
		  break
		}
	      case id := <-bodyDoneChan:
		if id == msgId {
		  //the body is been recived correctly, then do nothing
		  break
                }
	    }
    }
    
    msgId++
  }
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
		logger.Println("handle http req after creating tls connection:")
		logger.Println(req)
		req.Write(conn)
		startTimeoutChan <- true
		res, ok := http.ReadResponse(bufio.NewReader(conn), req)
		if ok != nil {
			logger.Println("error get Response")
		} else {
			logger.Println("get response correct")
		}
		reschan <- res
		//test!!!
		go func (){
			time.Sleep(10*time.Second)
			conn.Close()
		}()
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

	/*
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	server_conn, err = tls.Dial("tcp", url, conf)
	*/

	server_conn, err = net.Dial("tcp", url)
	if err != nil {
		logger.Println("error estabish connection with server\n")
		logger.Println(err)
		return nil
	}

	// so send the magic number to server, then server can identify
	// this is a connection for http not https
	_, err = server_conn.Write([]byte("YAEFCTqyz"+"\n")) // magic number for http connection
	if err != nil {
		logger.Println("jumpClient: error write")
		logger.Println(err)
		server_conn.Close()
		return nil
	}

	//do we need response from server, may not now

	return server_conn
}
