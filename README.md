# gojump
gojump jump from inner to outter

use jumpClient.go and jumpServer.go for jump.

jumpClient.go --compile--> JumpClient
jumpServer.go --compile--> JumpServer

Description:
  purpose of those two software is to jump Wall.

Application (IE or others) <--> JumpClient (as http proxy) <----- Wall ---> JumpServer <------> destination

prerequsite:
  JumpServer: need deploy this server on a outter hostt, so it can connect outter servers directly.
  JumpClient: jumpClinet need to connect with jumpServer directly.

Configurations: (client and server apps read the files(config files & tls files) which has hard-code name file names
  JumpClient need clientconfig.json 
  JumpServer need serverconfig.json
  JumpServer need server.crt and server.key for tls connection (you can generate yours)

  to run jumpClient:
  ./jumpClient

  to run jumpServer:
  ./jumpServer

Logs:
  JumpClient will generate jumpClient.log
  JumpServer will generate jumpServer.log


