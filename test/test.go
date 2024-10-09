package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/vmxy/go-ioi/ioi"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		ioi.Log.Println("========no args")
		//return
	}

	ioi.Log.Println("=====", args)
	host, port := "127.0.0.1", 9292
	if args[0] == "server" {
		ioi.Listen(host, port, func(session *ioi.Session) {
			ioi.Log.Println("accept==========")
			go HandleAccept(session)
			conn, err := session.OpenStream()
			if err != nil {
				session.Close()
				ioi.Log.Println("open stream error", err)
				return
			}
			conn.Write([]byte("server write client --->x1------------"))
			bs := make([]byte, 1024)
			size, err := conn.Read(bs)
			if err != nil {
				if err == io.EOF {
					return
				}
				ioi.Log.Println("read error 1", err)
				session.Close()
				return
			}
			fmt.Println("server open stream on read", string(bs[0:size]))
			conn.Close()
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
		go HandleAccept(session)
		conn, err := session.OpenStream()
		if err != nil {
			log.Panicln("dail error", err)
		}
		conn.Write([]byte("client write server ---> request"))
		bs := make([]byte, 1024)
		size, err := conn.Read(bs)
		if err != nil {
			if err == io.EOF {
				return
			}
			ioi.Log.Println("read error 2")
			return
		}
		fmt.Println("read by client=====>", string(bs[0:size]))
		conn.Close()
	}

}
func ParseInt(v string) int {
	x, e := strconv.ParseInt(v, 10, 32)
	if e != nil {
		return 0
	}
	return int(x)
}

func HandleAccept(session *ioi.Session) {
	defer session.Close()
	lis, err := session.Listen()
	if err != nil {
		ioi.Log.Println("listen error", lis)

	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			ioi.Log.Println("listen error", lis)
			break
		}
		go func() {
			bs := make([]byte, 1024)
			for {
				size, err := conn.Read(bs)
				if err != nil {
					conn.Close()
					break
				}
				fmt.Println("read accept data=", string(bs[0:size]))
				conn.Write([]byte("response--->" + string(bs[0:size])))
			}
		}()

	}

}
