package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	builtins "github.com/HoriaMercan/MapReducer/builtins"
)

type Student builtins.Student

func main() {
	s := &builtins.Student{10, "aaa", "vfd"}
	rpc.Register(s)
	rpc.HandleHTTP()
	sockname := "/var/tmp/echo.sock"
	os.Remove(sockname)
	listenObject, err := net.Listen("unix", sockname)
	if err != nil {
		log.Fatal("Listen error on server: ", err)
	}
	go http.Serve(listenObject, nil)

	for {
		rpc.Accept(listenObject)
	}
}
