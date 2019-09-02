package main
import (
    "crypto/tls"
    //"flag"
    "io"
    "log"
    "net"
    "net/http"
    "time"
    "httproad"
)

/************************************************************
hanlde outer tunnel. using httproad to send 
************************************************************/
func handleOuterTunnel(w http.ResponseWriter, r *http.Request) {
    // first need to tell jumpserver to create connection with
    //server    
    //set Authorization of Header to Id of the connection.
    id := httproad.GetConnId()
    //id to string
    r.Header.Set("Authorization",strconv.Itoa(id))
    
    //send http Connect request to jumpserver, if jumpserver
    //return okay, then we can build real http road
    err := httproad.SendConnectReq(r)
    if err != nil {
      httproad.FreeConnId(id)
      http.Error(w, err.Error(), http.StatusServiceUnavailable)
      return
    }
    
    w.WriteHeader(http.StatusOK)

    //hijack the unerlying connection for https transport
    hijacker, ok := w.(http.Hijacker)
    if !ok {
        http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
        return
    }
    client_conn, _, err := hijacker.Hijack()
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
    }
   
    //add this hijacked connection to httproad
    httproad.Addhttps(client_conn,id)
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
    // if it is outer request, then need handle outer
    if !httproad.IsInner(r){
      handleOuterTunnel(w,r)
      return
    }

    //handle inner tunnel
    dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
    hijacker, ok := w.(http.Hijacker)
    if !ok {
        http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
        return
    }
    client_conn, _, err := hijacker.Hijack()
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
    }
    go transfer(dest_conn, client_conn)
    go transfer(client_conn, dest_conn)
}
func transfer(destination io.WriteCloser, source io.ReadCloser) {
    defer destination.Close()
    defer source.Close()
    io.Copy(destination, source)
}


func handleHTTP(w http.ResponseWriter, req *http.Request) {
    var err error
    var resp http.Response
    // check inner host or outer host
    if httproad.IsInner(req){
      resp, err = http.DefaultTransport.RoundTrip(req)
    }else{
      //outer host, need use httproad
      resp, err = httproad.Sendhttp(req)
    }

    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }
    defer resp.Body.Close()
    copyHeader(w.Header(), resp.Header)
    w.WriteHeader(resp.StatusCode)
    io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
    for k, vv := range src {
        for _, v := range vv {
            dst.Add(k, v)
        }
    }
}
func main() {
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
