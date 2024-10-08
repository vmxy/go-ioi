package main

import (
	"flag"
	"fmt"

	"github.com/vmxy/go-ioi/ioi"
	"github.com/vmxy/go-ioi/ioi/stream"
	"github.com/vmxy/go-ioi/ioi/util"
)

func main() {

	fmt.Println("asdfasdf")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		util.Log.Println("========no args")
		//return
	}
	util.Log.Println("=====", args)
	host, port := "127.0.0.1", 9292
	if args[0] == "server" {
		ioi.Listen(host, port, func(session *stream.Session) {
			util.Log.Println("accept==========")
		})
	} else if args[0] == "client" {

	}

}
