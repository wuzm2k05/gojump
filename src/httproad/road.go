package httproad

/**********************************
This pacakge use Openshift docker service as road between goJump client and goJump
server.

This pacakge send http post to Openshift docker which is deployed with goJump server.

goJump client send http post as there is new pacakge in queue, and get response from goJump server, then distribute those pacakge in post response body to different tcp connections.

***************************************************/

type InnerMsg struct{
  type int // type of msg, 0: NULL req; 1: https req; 2: http req
  buf []byte  
}

msgChan := make(chan InnnerMsg, 100)
httpmsgChan := make(chan InnerMsg)


/*****************************************************************
check if this rquest will go to inner host or outer host
if inner host, that means not necessary go to jumpserver.
If it is a outer host, need jumpserver
*****************************************************************/
func IsInner(r *http.Request) bool {
  return false
}

/*****************************************************************
add one https hijacked tcp connection to httproad, road would distribute the pacakge to tcp connection, for those package in post response
*****************************************************************/
func Addhttps(conn net.Conn, connId int){
  connItem = getConnEntry(connId).conn
  connItem.conn = conn
  go singleConn(connItem)
}

/*********************************************************************** 
send one tcp pacaket to road. road would add it to queue and notify sender to send everthing in queue to goJump server 
***********************************************************************/
func sendpacket(){
}

/**********************************************************************
send one http packet to road. road would add it to queue and notify sender to send everthing i queue to goJump server. and wait until post response which include this http request response comes back, then return response to this function caller  
**********************************************************************/
func Sendhttp(r *http.Request){
  //serialize http request
  //send it to msgChan
  //wait response on httpmsgChan
}


