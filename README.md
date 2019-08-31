# gojump
gojump jump from inner to outter

use httpsClient.go and httpsServer.go for jump. other go files are not needed anymore

httpsClient.go --compile--> JumpClient
httpsServer.go --compile--> JumpServer

Description:
  purpose of those two software is to jump Wall.

Application (IE or others) <--> JumpClient (as http proxy) <----- Wall ---> JumpServer <------> destination

Configurations:
  JumpClient need clientconfig.json 
  JumpServer need serverconfig.json
  JumpServer need server.crt and server.key for tls connection (you can generate yours)

Logs:
  JumpClient will generate jumpClient.log
  JumpServer will generate jumpServer.log


