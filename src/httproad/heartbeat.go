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
        //send a null request to road
    }
  }
}
