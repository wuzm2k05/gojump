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
    }
    if er == EOF{
    }
    if er != nil {
    } 
  }
}

/***********************************************************************
send one tcp pacaket to road. road would add it to queue and notify sender to send everthing in queue to goJump server
***********************************************************************/
func sendpacket(){
}

