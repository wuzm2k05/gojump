package httproad

import (
  "time"
)

/***************************************
This file inlcudes a ticker which send a
 null post request to sender every 1 second
to pull those data from jumpserver
***************************************/
//bool hasdata
//int duration = 100

func hb_loop(){
  ticker := time.NewTicker(time.Second);

  //loop forever
  for {
    select {
      case <- ticker.C:
        sendNullMsg()
    }
  }
}

/************************************************************************
send one packet which indicate the NULL request.
*************************************************************************/
func sendNullMsg(){
  msg := innerMsg{0,0}
  msgChan <- &msg
end one packet which indicate the NULL request.
*************************************************************************/
func sendNullMsg(){
}


