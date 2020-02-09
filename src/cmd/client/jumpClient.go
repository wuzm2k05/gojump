package main

import (
	"crypto/tls"
	//"io/ioutil"
	"encoding/json"
	"fmt"
	"httproad"
	"io"
	"log"
	"net/http"
	"os"
	"slog"
	"strings"
	//"time"
)

type ClientConf struct {
	JumpClientListenAddr, JumpServerUrl string
}

var logger *log.Logger
var clientConf *ClientConf

func parseClientConf() {
	file, err := os.Open("clientconf.json")
	if err != nil {
		panic("no clientconf.json file, need configurationfile!!")
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	clientConf = &ClientConf{JumpClientListenAddr: ":8888"}
	err = decoder.Decode(&clientConf)
	if err != nil {
		panic("decode clientconf file fail!!")
	}
}

/**
Des:
  need to create tls with jump server.
  Then send connect request to jump server.
  then hijack the connection.
**/
func handleTunneling(w http.ResponseWriter, r *http.Request) {

	logger.Println("get a Connect request from app\n")

	//create a tls connection with jump server
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	server_conn, err := tls.Dial("tcp", clientConf.JumpServerUrl, conf)
	if err != nil {
		logger.Println("error estabish connection with server\n")
		logger.Println(err)
		return
	}

	//send destination host:port to jump server,
	//so the jump server can connect with destination host:port
	n, err := server_conn.Write([]byte("XAEFCTqyz " + r.Host + "\n")) // magic number + host
	if err != nil {
		logger.Println("jumpClient: error write")
		logger.Println(err)
		server_conn.Close()
		return
	}

	//wait server response, if server response okay,
	//then send okay to app; otherwise send faile to app
	buf := make([]byte, 16)
	n, err = server_conn.Read(buf)
	if err != nil {
		logger.Println("jumpClient: error read")
		logger.Println(err)
		server_conn.Close()
		return
	}
	if "okay" != string(buf[:n]) {
		logger.Println("error return code from server:" + string(buf[:n]))
		w.WriteHeader(http.StatusForbidden)
		server_conn.Close()
		return
	}

	w.WriteHeader(http.StatusOK)

	//hijack client connection
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		logger.Println("hijiack not support!\n")
		server_conn.Close()
		return
	}
	client_conn, _, _ := hijacker.Hijack()
	go transfer(server_conn, client_conn)
	go transfer(client_conn, server_conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

/*
is the buffer have end of html </html>
we only support text/html content type now
*/
func isBodyDone(buf []byte) bool{
  length := len(buf)
  if length > 100 {
    if strings.Contains(string(buf[length-100:]),"</html>") {
      return true
    } else {
      return false
    }
  } else {
    if strings.Contains(string(buf),"</html>"){
      return true
    }else{
      return false
    }
  }
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf(w, "Hi, http protocl is NOT support, you may want use https instead!")
	//return // not support HTTP now	
	//defer res.Body.Close()
	res,msgId := httproad.SendHttpReq(clientConf.JumpServerUrl, req)


	if res.StatusCode != http.StatusOK {
		httproad.RecResBodyDone(msgId)
		logger.Println("http response is not OK")
		fmt.Fprintf(w,"error happen: statuscode: "+string(res.StatusCode))
	}

	//res.Body.Read only read the current arrived content. However, some body content has not arrived yet!!
	//even the response is got, that only means the header of response arrived, but not all body!!
	idx := 0
	var buf [4194304]byte
	for {
		n, err := res.Body.Read(buf[idx:])
		if err != nil {
			logger.Println("read body error!")
			logger.Println(err)
		}
		idx += n
		if isBodyDone(buf[:idx]){
			break
		}
		logger.Println("read length: ")
		logger.Println(n)
	}
	
        //notify http road the body is done
        httproad.RecResBodyDone(msgId)
	//bodySS := string(buf[:])
	fmt.Fprintf(w,string(buf[:idx]))

	return
}

func main() {
	slog.LoggerInit("jumpClient.log")
	logger = slog.GetInstance()
	logger.Println("jumpClient start")

	parseClientConf()

	//start http proxy server for app connect
	server := &http.Server{
		Addr: clientConf.JumpClientListenAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				handleTunneling(w, r)
			} else {
				handleHTTP(w, r)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	logger.Println("JumpClient listen on Addr:" + clientConf.JumpClientListenAddr)
	log.Fatal(server.ListenAndServe())
}