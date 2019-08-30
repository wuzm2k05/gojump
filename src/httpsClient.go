package main

import (
  "crypto/tls"
//  "io/ioutil"
  "fmt"
  "net/http"
  "log"
  "io"
  "slog"
)

var logger *log.Logger

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

  server_conn, err := tls.Dial("tcp", "127.0.0.1:4444",conf)
  if err != nil {
    logger.Println("error estabish connection with server\n")
    logger.Println(err)
    return
  }

  //send destination host:port to jump server,
  //so the jump server can connect with destination host:port
  n, err := server_conn.Write([]byte("XAEFCTqyz "+r.Host+"\n"))  // magic number + host
  if err != nil {
    logger.Println("jumpClient: error write")
    logger.Println(err)
    server_conn.Close()
    return
  }

  //wait server response, if server response okay,
  //then send okay to app; otherwise send faile to app
  buf := make([]byte,16)
  n, err = server_conn.Read(buf)
  if err != nil{
    logger.Println("jumpClient: error read")
    logger.Println(err)
    server_conn.Close()
    return
  }
  if "okay" != string(buf[:n]){
    logger.Println("error return code from server:"+string(buf[:n]))
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
  client_conn, _, _:= hijacker.Hijack()
  go transfer(server_conn, client_conn)
  go transfer(client_conn, server_conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
    defer destination.Close()
    defer source.Close()
    io.Copy(destination, source)
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(w, "Hi, NOT support!")
}

func main(){
 
  slog.LoggerInit("jumpClient.log")
  logger = slog.GetInstance()
  logger.Println("jumpClient start")
  logger.Println("jumpClient start")
 
  //start http proxy server for app connect
  server := &http.Server{
    Addr: ":8888",
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
  log.Fatal(server.ListenAndServe())

}
