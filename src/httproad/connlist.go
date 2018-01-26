package httproad

type connList struct{
  //id int
  closing bool  //peer closed
  isUsed bool   //this connection id is used 
  conn net.Conn //connection of tcp
}

var httpsConnList [100]connList

/*****************************************************
get one connectionId from list
[API] means it is visible outside of pacakge
*****************************************************/
func GetConnId() int{
  for i:=0; i<100; i++{
    if (!httpsConnList[i].isUsed){
      httpsConnList[i].isUsed = true
      httpsConnList[i].closing = false
      return i
    }
  }
  return 0
}

/****************************************************
free one connection Id
[API]
*****************************************************/
func FreeConnId(id int){
  if id<100{
    httpsConnList[id].isUsed = false
    httpsConnList[id].closing = false
    httpsConnList[id].conn = nil
  }
}


/****************************************************
get one entry of connection list
****************************************************/
func getConnEntry(id int) *connList{
  if id<100{
    return &(httpsConnList[id])
  }
  return nil  
}

/**************************************************
closing one connection. Peer already closed the connection.
so need this site to close the connection too
**************************************************/
func closConn(id int){
 //may don't need this func 
}
