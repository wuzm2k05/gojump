package main

import (
  "crypto/tls"
//  "io/ioutil"
  "fmt"
  "net/http"
  "log"
  "io"
)

/**
Des:
  create a tls connection with jump server
Return value:
  net connection and error information
**/
func createTLSwithServer() {
  //use local jump server for testing
  
}

/**
Des:
  need to create tls with jump server. 
  Then send connect request to jump server.
  then hijack the connection.
**/
func handleTunneling(w http.ResponseWriter, r *http.Request) {

  //create a tls connection with jump server
  conf := &tls.Config{
    InsecureSkipVerify: true,
  }

  server_conn, err := tls.Dial("tcp", "127.0.0.1:4444",conf)
  if err != nil {
    fmt.Println("error estabish connection with server\n")
    fmt.Println(err)
    return
  }

  //send destination host:port to jump server,
  //so the jump server can connect with destination host:port
  fmt.Println("connect to Host:"+r.Host)
  n, err := server_conn.Write([]byte("XAEFCTqyz "+r.Host+"\n"))  // magic number + host
  if err != nil {
    fmt.Println("error write")
    return
  }

  //wait server response, if server response okay,
  //then send okay to app; otherwise send faile to app
  buf := make([]byte,16)
  n, err = server_conn.Read(buf)
  if err != nil{
    fmt.Println("error read")
    return
  }
  fmt.Println("recieve from server:"+string(buf[:n]))
  if "okay" != string(buf[:n]){
    fmt.Println("error return code from server:"+string(buf[:n]))
    w.WriteHeader(http.StatusForbidden)
    return
  }

  w.WriteHeader(http.StatusOK)

  //hijack client connection 
  hijacker, ok := w.(http.Hijacker) 
  if !ok {
    fmt.Println("hijiack not support!\n")
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

/*
  tr := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
  }
  client := &http.Client{Transport: tr}
  resp, err := client.Get("https://127.0.0.1:9090")

  if err != nil {
    fmt.Println("error:", err)
    return
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  fmt.Println(string(body))
*/
}
