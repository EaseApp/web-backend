package main

import(
  "bytes"
  "net"
  "log"
)


func handleConnection(conn net.Conn){
  // defer conn.Close()

  buf := new(bytes.Buffer)
  buf.ReadFrom(conn)
  s := buf.String() // Does a complete copy of the bytes in the buffer.
  log.Println(&conn, s)
}

func main(){
  ln, err := net.Listen("tcp", ":8081")
  if err != nil {
  	log.Println(err)
  }
  log.Println("Listening on port 8081")
  for {
  	conn, err := ln.Accept()
  	if err != nil {
  		log.Println(err)
  	}
  	go handleConnection(conn)
  }
}
