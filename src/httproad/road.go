package httproad

/**********************************
This pacakge use Openshift docker service as road between goJump client and goJump
server.

This pacakge send http post to Openshift docker which is deployed with goJump server.

goJump client send http post as there is new pacakge in queue, and get response from goJump server, then distribute those pacakge in post response body to different tcp connections.

***************************************************/

type innerMsg struct{
  // type of msg, 0: heartbeat req; 1: https req; 2: http connect req;
  //              3: connection close;
  type int 
  connId int//connection ID if it is a https req
  buf []byte  // mesg content
}

msgChan := make(chan *InnnerMsg, 100)
httpmsgChan := make(chan *innerMsg)


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

/**********************************************************************
send one http packet to road. road would add it to queue and notify sender to send everthing i queue to goJump server. and wait until post response which include this http request response comes back, then return response to this function caller  

!!!NOT support http right now.
**********************************************************************/
func Sendhttp(r *http.Request){
  //serialize http request
  //send it to msgChan
  //wait response on httpmsgChan
}

/*************************************************************************
send one http connect request to goJump server. wait until post response.
encode the connectReq and send it through httpChan
*************************************************************************/
func SendConnectReq(r *http.Request,id int)(error){
  msg := innerMsg{2,id}
  msg.bug = []byte(r.Host)
  httpmsgChan <- &msg
  res := <- httpmsgChan
  
}
