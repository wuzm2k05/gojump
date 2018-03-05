package httproad

import (
  "os"
  "bufio"
  "strings"
)

/*******************************************************************
this function would be the core func to send the packet to server
*******************************************************************/

func longroad(){
  //wait msg on msgChan and httpmsgChan
  for {
    select {
    case msg := <-msgChan:
      handleMsgFromJumpServer(sendMsgToJumpServer(msg))
    case msg := <-httpmsgChan:
      httpmsgChan <- sendMsgToJumpServer(msg)
    }
  }
}

/****************************************************************
send msg to jump server. 
1> The msg will be encoded as post http req body.
2> The http post req will be send to JumpServer.
3> Then waiting for jumpServer to return the http post req response,
4> Then return the response to caller.
*****************************************************************/
func sendMsgToJumpServer(msg innerMsg){

  // get url from local setting file
  url := ""
  f,_ := os.Open("goJump.cfg")
  if f != nil{
    buf := bufio.NewReader(f)
    line, err := buf.ReadString('\n')
    if err != nil {
      fmt.Printf("Error: Open config file goJump.cfg fail!")
    }
    url = strings.TrimSpace(line)
  }
  
  //seriiza innerMsg
  ser_inner := ""

  //send the msg to jump server, and waiting for response
  resp, err := http.Post(url, "application/x-www-form-urlencoded", ser_inner)
  if err != nil {
    fmt.Println(err)
  }

  return resp, err
}

/*****************************************************************
handle the http post req response from JumpServer.
The response could include more than one msg. so need to parse all 
of them.
*****************************************************************/
func handleMsgFromJumpServer(){
}

/******************************************************************
singleConn: this is the goroutine to recieve pacaket from application.
******************************************************************/
func singleConn(connItem *connList){
  //wait data and send data
  src := connItem.conn
  buf := make([]byte,32*1024)
  for {
    nr, er := src.Read(buf)
    if nr > 0{
      sendPacket(buf[0:nr],connList.connId)
    }
    if er == EOF{
      sendClosePacket(connList.connId)
      break;
    }
    if er != nil {
      sendClosePacket(connList.connId)
      break;
    } 
  }
}

/***********************************************************************
send one tcp pacaket to road. road would add it to queue and notify sender to send everthing in queue to goJump server
***********************************************************************/
func sendPacket(buf []byte, connId int) {
  msg := innerMsg{1,connId,buf}
  msgChan <- &msg
}

/*************************************************************************
send one packet which indicate the connection is closed.
**************************************************************************/
func sendClosePacket(connId int){
  msg := innerMsg{3,connId}
  msgChan <- &msg
}
