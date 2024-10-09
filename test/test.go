package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/vmxy/go-ioi/ioi"
	"github.com/vmxy/go-ioi/ioi/stream"
	"github.com/vmxy/go-ioi/ioi/util"
)

func main() {
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
		var f flag.FlagSet
		f.Usage = func() {
			f.PrintDefaults()
		}
		server := f.String("server", "127.0.0.1:9393", "server host:port")
		f.Parse(args[1:])

		fmt.Println("server", *server, args)

		hp := strings.Split(*server, ":")

		host, port := hp[0], ParseInt(hp[1])

		session, err := ioi.Dail(host, port)
		if err != nil {
			log.Panicln("dail error", err)
		}
		conn, err := session.OpenStream()
		if err != nil {
			log.Panicln("dail error", err)
		}
		conn.Write([]byte("request"))
	}

}
func ParseInt(v string) int {
	x, e := strconv.ParseInt(v, 10, 32)
	if e != nil {
		return 0
	}
	return int(x)
}
