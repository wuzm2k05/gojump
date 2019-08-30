package slog

import (
  "log"
  "os"
  "fmt"
)

var instance *log.Logger

func LoggerInit(name string){
  file, err := os.Create(name)
  if err != nil {
    fmt.Println("fail to create "+name+" log file");
    return
  }
  instance =log.New(file, "", log.LstdFlags|log.Llongfile) 
}

func GetInstance() *log.Logger {
  return instance
}


